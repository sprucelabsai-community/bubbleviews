package render

import (
	"strings"

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
	case child.Columns != nil:
		return renderColumns(*child.Columns, parentSize)
	case child.List != nil:
		return renderList(*child.List, parentSize)
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

	interiorWidth := 0
	if frame.FillWidth && parentSize.Width > 0 {
		interiorWidth = max(parentSize.Width-2, 0)
	} else if frame.Content != nil && frame.Content.Size.Width > 0 {
		interiorWidth = frame.Content.Size.Width + frame.Padding.Left + frame.Padding.Right
	}
	if interiorWidth > 0 {
		style = style.Width(interiorWidth)
	}

	interiorHeight := 0
	if frame.FillHeight && parentSize.Height > 0 {
		interiorHeight = max(parentSize.Height-2, 0)
	} else if frame.Content != nil && frame.Content.Size.Height > 0 {
		interiorHeight = frame.Content.Size.Height + frame.Padding.Top + frame.Padding.Bottom
	}
	if interiorHeight > 0 {
		style = style.Height(interiorHeight)
	}

	if frame.Content == nil {
		return style.Render("")
	}

	contentWidth := 0
	if frame.Content.Size.Width > 0 {
		contentWidth = frame.Content.Size.Width
	} else if interiorWidth > 0 {
		contentWidth = max(interiorWidth-frame.Padding.Left-frame.Padding.Right, 0)
	}

	contentHeight := 0
	if frame.Content.Size.Height > 0 {
		contentHeight = frame.Content.Size.Height
	} else if interiorHeight > 0 {
		contentHeight = max(interiorHeight-frame.Padding.Top-frame.Padding.Bottom, 0)
	}

	contentView := *frame.Content
	contentView.Size = bubbleviews.Size{
		Width:  contentWidth,
		Height: contentHeight,
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

func renderColumns(columns bubbleviews.ColumnsView, parentSize bubbleviews.Size) string {
	if len(columns.Columns) == 0 {
		return ""
	}

	widths := computeColumnWidths(columns, parentSize.Width)

	rendered := make([]string, len(columns.Columns))
	maxHeight := 0

	for i, col := range columns.Columns {
		if col.Content == nil {
			rendered[i] = lipgloss.NewStyle().Width(widths[i]).Render("")
		} else {
			childView := *col.Content
			childView.Size = bubbleviews.Size{
				Width:  widths[i],
				Height: parentSize.Height,
			}

			rendered[i] = renderView(childView)
		}

		if h := lipgloss.Height(rendered[i]); h > maxHeight {
			maxHeight = h
		}
	}

	segments := make([]string, 0, len(columns.Columns)*2-1)
	for i, segment := range rendered {
		if i > 0 && columns.Spacing > 0 {
			segments = append(segments, lipgloss.NewStyle().Width(columns.Spacing).Render(""))
		}

		segments = append(segments, segment)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, segments...)
}

func computeColumnWidths(columns bubbleviews.ColumnsView, parentWidth int) []int {
	count := len(columns.Columns)
	widths := make([]int, count)

	if count == 0 {
		return widths
	}

	available := parentWidth
	if available > 0 && count > 1 {
		available -= columns.Spacing * (count - 1)
		if available < 0 {
			available = 0
		}
	}

	flexible := make([]int, 0, count)

	for i, col := range columns.Columns {
		if col.Width > 0 {
			if available > 0 {
				if col.Width > available {
					widths[i] = available
				} else {
					widths[i] = col.Width
				}
				available -= widths[i]
			} else {
				widths[i] = col.Width
			}
		} else {
			flexible = append(flexible, i)
		}
	}

	if available < 0 {
		available = 0
	}

	share := 0
	if len(flexible) > 0 && available > 0 {
		share = available / len(flexible)
	}

	for _, idx := range flexible {
		widths[idx] = share
	}

	remainder := 0
	if len(flexible) > 0 {
		remainder = available - share*len(flexible)
		for i := 0; i < remainder && i < len(flexible); i++ {
			widths[flexible[i]]++
		}
	}

	return widths
}

func renderList(list bubbleviews.ListView, parentSize bubbleviews.Size) string {
	lines := make([]string, 0, len(list.Items)+1)

	availableWidth := parentSize.Width

	if list.Title != "" {
		style := lipgloss.NewStyle().Bold(true)
		if color := string(list.TitleColor); color != "" {
			style = style.Foreground(lipgloss.Color(color))
		}

		if availableWidth > 0 {
			for _, wrapped := range wrapText(list.Title, availableWidth) {
				lines = append(lines, style.Render(wrapped))
			}
		} else {
			lines = append(lines, style.Render(list.Title))
		}
	}

	itemStyle := lipgloss.NewStyle()
	if color := string(list.ItemColor); color != "" {
		itemStyle = itemStyle.Foreground(lipgloss.Color(color))
	}

	bullet := list.Bullet
	if bullet == "" {
		bullet = "- "
	}
	bulletWidth := lipgloss.Width(bullet)
	textWidth := availableWidth - bulletWidth
	if availableWidth <= 0 {
		textWidth = -1
	}
	if textWidth < 1 && availableWidth > 0 {
		textWidth = 1
	}
	bulletPadding := strings.Repeat(" ", max(bulletWidth, 0))

	for _, item := range list.Items {
		wrapped := []string{item}
		if textWidth > 0 {
			wrapped = wrapText(item, textWidth)
		}

		for i, segment := range wrapped {
			prefix := bullet
			if i > 0 {
				prefix = bulletPadding
			}
			lines = append(lines, itemStyle.Render(prefix+segment))
		}
	}

	if len(lines) == 0 {
		return ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0)
	current := words[0]

	for _, word := range words[1:] {
		if lipgloss.Width(current+" "+word) <= width {
			current += " " + word
			continue
		}

		lines = append(lines, current)
		current = word
	}

	lines = append(lines, current)

	return lines
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
