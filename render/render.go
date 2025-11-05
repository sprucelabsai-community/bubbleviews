package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/sprucelabsai-community/bubbleviews"
)

// Render converts a View tree into a fully formatted string.
func Render(view bubbleviews.View) string {
	return renderView(view)
}

func renderView(view bubbleviews.View) string {
	if len(view.Children) == 0 {
		return ""
	}

	outputs := make([]string, 0, len(view.Children))
	for _, child := range view.Children {
		if rendered := renderChild(child, view.Size); rendered != "" {
			outputs = append(outputs, rendered)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, outputs...)
}

func renderChild(child bubbleviews.ViewChild, parentSize bubbleviews.Size) string {
	switch {
	case child.Frame != nil:
		return renderFrame(*child.Frame, parentSize)
	case child.Placement != nil:
		return renderPlacement(*child.Placement, parentSize)
	case child.Button != nil:
		return renderButton(*child.Button)
	default:
		return ""
	}
}

func renderFrame(frame bubbleviews.FrameView, parentSize bubbleviews.Size) string {
	style := lipgloss.NewStyle()

	if border := mapBorderStyle(frame.Border); border != nil {
		style = style.BorderStyle(*border)
	}

	if color := string(frame.BorderColor); color != "" {
		style = style.BorderForeground(lipgloss.Color(color))
	}

	style = style.Padding(frame.Padding.Top, frame.Padding.Right, frame.Padding.Bottom, frame.Padding.Left)

	width := parentSize.Width
	if !frame.FillWidth && frame.Content != nil && frame.Content.Size.Width > 0 {
		width = frame.Content.Size.Width + 2 + frame.Padding.Left + frame.Padding.Right
	}
	if width > 0 {
		style = style.Width(width)
	}

	height := parentSize.Height
	if !frame.FillHeight && frame.Content != nil && frame.Content.Size.Height > 0 {
		height = frame.Content.Size.Height + 2 + frame.Padding.Top + frame.Padding.Bottom
	}
	if height > 0 {
		style = style.Height(height)
	}

	if frame.Content == nil {
		return style.Render("")
	}

	contentView := *frame.Content
	contentView.Size = bubbleviews.Size{
		Width:  max(width-2-frame.Padding.Left-frame.Padding.Right, 0),
		Height: max(height-2-frame.Padding.Top-frame.Padding.Bottom, 0),
	}

	return style.Render(renderView(contentView))
}

func renderPlacement(placement bubbleviews.PlacementView, parentSize bubbleviews.Size) string {
	if placement.Content == nil {
		return ""
	}

	areaWidth := placement.AreaWidth
	areaHeight := placement.AreaHeight

	if areaWidth == 0 {
		areaWidth = parentSize.Width
	}
	if areaHeight == 0 {
		areaHeight = parentSize.Height
	}

	contentView := *placement.Content
	contentView.Size = bubbleviews.Size{
		Width:  max(areaWidth, 0),
		Height: max(areaHeight, 0),
	}

	child := renderView(contentView)

	return lipgloss.Place(
		max(areaWidth, lipgloss.Width(child)),
		max(areaHeight, lipgloss.Height(child)),
		mapHorizontal(placement.HAlign),
		mapVertical(placement.VAlign),
		child,
	)
}

func renderButton(button bubbleviews.ButtonView) string {
	style := lipgloss.NewStyle()

	if border := mapBorderStyle(button.Border); border != nil {
		style = style.BorderStyle(*border)
	}

	if color := string(button.BorderColor); color != "" {
		style = style.BorderForeground(lipgloss.Color(color))
	}

	if text := string(button.TextColor); text != "" {
		style = style.Foreground(lipgloss.Color(text))
	}

	style = style.Padding(button.Padding.Top, button.Padding.Right, button.Padding.Bottom, button.Padding.Left)

	return style.Render(button.Label)
}

func mapBorderStyle(style bubbleviews.BorderStyle) *lipgloss.Border {
	switch style {
	case bubbleviews.BorderThin:
		border := lipgloss.NormalBorder()
		return &border
	case bubbleviews.BorderThick:
		border := lipgloss.ThickBorder()
		return &border
	case bubbleviews.BorderNone:
		return nil
	default:
		return nil
	}
}

func mapHorizontal(align bubbleviews.Alignment) lipgloss.Position {
	switch align {
	case bubbleviews.AlignCenter:
		return lipgloss.Center
	case bubbleviews.AlignEnd:
		return lipgloss.Right
	default:
		return lipgloss.Left
	}
}

func mapVertical(align bubbleviews.Alignment) lipgloss.Position {
	switch align {
	case bubbleviews.AlignCenter:
		return lipgloss.Center
	case bubbleviews.AlignEnd:
		return lipgloss.Bottom
	default:
		return lipgloss.Top
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
