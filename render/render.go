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
		if rendered := renderNode(child, view.Size); rendered != "" {
			outputs = append(outputs, rendered)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, outputs...)
}

func renderNode(node bubbleviews.Node, parentSize bubbleviews.Size) string {
	switch n := node.(type) {
	case bubbleviews.BoxNode:
		return renderBox(n, parentSize)
	case *bubbleviews.BoxNode:
		return renderBox(*n, parentSize)
	case bubbleviews.FlexNode:
		return renderFlex(n, parentSize)
	case *bubbleviews.FlexNode:
		return renderFlex(*n, parentSize)
	case bubbleviews.FlowNode:
		return renderFlow(n, parentSize)
	case *bubbleviews.FlowNode:
		return renderFlow(*n, parentSize)
	case bubbleviews.ASCIIArtNode:
		return renderASCIIArt(n, parentSize)
	case *bubbleviews.ASCIIArtNode:
		return renderASCIIArt(*n, parentSize)
	case bubbleviews.TextNode:
		return renderText(n, parentSize)
	case *bubbleviews.TextNode:
		return renderText(*n, parentSize)
	default:
		return ""
	}
}

func renderBox(box bubbleviews.BoxNode, parentSize bubbleviews.Size) string {
	style := lipgloss.NewStyle()

	if border := mapBorderStyle(box.Style.Border); border != nil {
		style = style.BorderStyle(*border)
	}

	if color := string(box.Style.BorderColor); color != "" {
		style = style.BorderForeground(lipgloss.Color(color))
	}

	style = style.Padding(
		box.Style.Padding.Top,
		box.Style.Padding.Right,
		box.Style.Padding.Bottom,
		box.Style.Padding.Left,
	)

	contentWidth := box.Content.Size.Width
	if box.Style.FillWidth && parentSize.Width > 0 {
		contentWidth = max(parentSize.Width-box.Style.Padding.Left-box.Style.Padding.Right-2, 0)
	}
	if contentWidth < 0 {
		contentWidth = 0
	}

	contentHeight := box.Content.Size.Height
	if box.Style.FillHeight && parentSize.Height > 0 {
		contentHeight = max(parentSize.Height-box.Style.Padding.Top-box.Style.Padding.Bottom-2, 0)
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	totalWidth := 0
	if contentWidth > 0 {
		totalWidth = contentWidth + box.Style.Padding.Left + box.Style.Padding.Right
	}
	if box.Style.FillWidth && parentSize.Width > 0 {
		totalWidth = max(parentSize.Width-2, 0)
	}
	if totalWidth > 0 {
		style = style.Width(totalWidth)
	}

	totalHeight := 0
	if contentHeight > 0 {
		totalHeight = contentHeight + box.Style.Padding.Top + box.Style.Padding.Bottom
	}
	if box.Style.FillHeight && parentSize.Height > 0 {
		totalHeight = max(parentSize.Height-2, 0)
	}
	if totalHeight > 0 {
		style = style.Height(totalHeight)
	}

	if len(box.Content.Children) == 0 {
		return style.Render("")
	}

	contentView := box.Content
	contentView.Size = bubbleviews.Size{
		Width:  contentWidth,
		Height: contentHeight,
	}

	contentRendered := renderView(contentView)
	if contentRendered == "" {
		return style.Render("")
	}

	widthForAlign := contentWidth
	if widthForAlign == 0 {
		widthForAlign = lipgloss.Width(contentRendered)
	}
	heightForAlign := contentHeight
	if heightForAlign == 0 {
		heightForAlign = lipgloss.Height(contentRendered)
	}

	if box.Style.HAlign != bubbleviews.AlignStart || box.Style.VAlign != bubbleviews.AlignStart {
		contentRendered = lipgloss.Place(
			max(widthForAlign, lipgloss.Width(contentRendered)),
			max(heightForAlign, lipgloss.Height(contentRendered)),
			mapHorizontal(box.Style.HAlign),
			mapVertical(box.Style.VAlign),
			contentRendered,
		)
	} else if widthForAlign > 0 {
		contentRendered = lipgloss.Place(
			widthForAlign,
			lipgloss.Height(contentRendered),
			lipgloss.Left,
			lipgloss.Top,
			contentRendered,
		)
	}

	return style.Render(contentRendered)
}

func renderFlex(flex bubbleviews.FlexNode, parentSize bubbleviews.Size) string {
	if len(flex.Items) == 0 {
		return ""
	}

	switch flex.Direction {
	case bubbleviews.FlexDirectionColumn:
		return renderFlexColumn(flex, parentSize)
	default:
		return renderFlexRow(flex, parentSize)
	}
}

func renderFlexRow(flex bubbleviews.FlexNode, parentSize bubbleviews.Size) string {
	widths := computeFlexWidths(flex, parentSize.Width)
	rendered := make([]string, len(flex.Items))
	maxHeight := 0

	for i, item := range flex.Items {
		childSize := bubbleviews.Size{
			Width:  widths[i],
			Height: parentSize.Height,
		}
		rendered[i] = renderNode(item.Node, childSize)

		if widths[i] > 0 {
			rendered[i] = lipgloss.Place(
				widths[i],
				lipgloss.Height(rendered[i]),
				lipgloss.Left,
				lipgloss.Top,
				rendered[i],
			)
		}

		if h := lipgloss.Height(rendered[i]); h > maxHeight {
			maxHeight = h
		}
	}

	segments := make([]string, 0, len(flex.Items)*2-1)
	for i, segment := range rendered {
		if i > 0 && flex.Spacing > 0 {
			segments = append(segments, lipgloss.NewStyle().Width(flex.Spacing).Render(""))
		}
		segments = append(segments, segment)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, segments...)
}

func renderFlexColumn(flex bubbleviews.FlexNode, parentSize bubbleviews.Size) string {
	segments := make([]string, 0, len(flex.Items)*2-1)

	for i, item := range flex.Items {
		childSize := bubbleviews.Size{
			Width:  parentSize.Width,
			Height: item.Height,
		}

		segment := renderNode(item.Node, childSize)
		if parentSize.Width > 0 {
			segment = lipgloss.Place(
				parentSize.Width,
				lipgloss.Height(segment),
				lipgloss.Left,
				lipgloss.Top,
				segment,
			)
		}

		if i > 0 && flex.Spacing > 0 {
			spacer := lipgloss.NewStyle().Height(flex.Spacing).Render("")
			segments = append(segments, spacer)
		}

		segments = append(segments, segment)
	}

	return lipgloss.JoinVertical(lipgloss.Left, segments...)
}

func renderFlow(flow bubbleviews.FlowNode, parentSize bubbleviews.Size) string {
	if len(flow.Items) == 0 {
		return ""
	}

	itemSpacing := flow.ItemSpacing
	if itemSpacing < 0 {
		itemSpacing = 0
	}
	rowSpacing := flow.RowSpacing
	if rowSpacing < 0 {
		rowSpacing = 0
	}

	minWidth := flow.ItemMinWidth
	if minWidth < 1 {
		minWidth = 1
	}

	maxColumns := len(flow.Items)
	if parentSize.Width > 0 {
		available := parentSize.Width + itemSpacing
		denominator := minWidth + itemSpacing
		if denominator > 0 {
			columns := available / denominator
			if columns < 1 {
				columns = 1
			}
			if columns < maxColumns {
				maxColumns = columns
			}
		}
	}
	if maxColumns < 1 {
		maxColumns = 1
	}

	var columnWidth int
	if parentSize.Width > 0 {
		columnWidth = parentSize.Width - itemSpacing*(maxColumns-1)
		if columnWidth < 0 {
			columnWidth = 0
		}
		if maxColumns > 0 {
			columnWidth = columnWidth / maxColumns
		}
	}

	rows := make([]string, 0, (len(flow.Items)+maxColumns-1)/maxColumns)

	for start := 0; start < len(flow.Items); start += maxColumns {
		end := start + maxColumns
		if end > len(flow.Items) {
			end = len(flow.Items)
		}

		segments := make([]string, 0, (end-start)*2-1)
		for i := start; i < end; i++ {
			size := bubbleviews.Size{Width: columnWidth, Height: parentSize.Height}
			segment := renderNode(flow.Items[i], size)
			segments = append(segments, segment)
		}

		row := joinWithSpacing(segments, itemSpacing, lipgloss.JoinHorizontal, lipgloss.Top)
		rows = append(rows, row)
	}

	return joinWithSpacing(rows, rowSpacing, lipgloss.JoinVertical, lipgloss.Left)
}

func renderASCIIArt(art bubbleviews.ASCIIArtNode, parentSize bubbleviews.Size) string {
	if len(art.Lines) == 0 {
		return ""
	}

	style := lipgloss.NewStyle()
	if color := string(art.Color); color != "" {
		style = style.Foreground(lipgloss.Color(color))
	}
	if art.Bold {
		style = style.Bold(true)
	}

	position := mapHorizontal(art.Align)
	if art.Align == "" {
		position = lipgloss.Left
	}

	styled := make([]string, len(art.Lines))
	for i, line := range art.Lines {
		rendered := style.Render(line)
		styled[i] = rendered
	}

	block := strings.Join(styled, "\n")
	if parentSize.Width > 0 {
		return lipgloss.PlaceHorizontal(parentSize.Width, position, block)
	}

	return block
}

func computeFlexWidths(flex bubbleviews.FlexNode, parentWidth int) []int {
	count := len(flex.Items)
	widths := make([]int, count)
	if count == 0 {
		return widths
	}

	available := parentWidth
	if available > 0 && count > 1 {
		available -= flex.Spacing * (count - 1)
		if available < 0 {
			available = 0
		}
	}

	flexible := make([]int, 0, count)
	growWeights := make([]int, count)
	totalGrow := 0

	for i, item := range flex.Items {
		if item.Width > 0 {
			if available > 0 {
				if item.Width > available {
					widths[i] = available
				} else {
					widths[i] = item.Width
				}
				available -= widths[i]
			} else {
				widths[i] = item.Width
			}
		} else {
			flexible = append(flexible, i)
			if item.Grow > 0 {
				growWeights[i] = item.Grow
				totalGrow += item.Grow
			}
		}
	}

	if available < 0 {
		available = 0
	}

	if len(flexible) > 0 {
		if totalGrow <= 0 {
			share := 0
			if available > 0 {
				share = available / len(flexible)
			}
			for _, idx := range flexible {
				widths[idx] = share
			}
			remainder := available - share*len(flexible)
			for i := 0; i < remainder && i < len(flexible); i++ {
				widths[flexible[i]]++
			}
		} else {
			for _, idx := range flexible {
				grow := growWeights[idx]
				width := 0
				if grow > 0 && available > 0 {
					width = available * grow / totalGrow
				}
				widths[idx] = width
			}
			assigned := 0
			for _, idx := range flexible {
				assigned += widths[idx]
			}
			remainder := available - assigned
			if remainder < 0 {
				remainder = 0
			}
			for i := 0; i < remainder && i < len(flexible); i++ {
				widths[flexible[i]]++
			}
		}
	}

	return widths
}

func renderText(text bubbleviews.TextNode, parentSize bubbleviews.Size) string {
	style := lipgloss.NewStyle()

	if color := string(text.Color); color != "" {
		style = style.Foreground(lipgloss.Color(color))
	}

	if text.Bold {
		style = style.Bold(true)
	}

	width := parentSize.Width
	prefix := text.Prefix
	continuation := text.ContinuationPrefix
	if continuation == "" {
		continuation = strings.Repeat(" ", lipgloss.Width(prefix))
	}

	wrapWidth := width
	if width > 0 {
		prefixWidth := lipgloss.Width(prefix)
		wrapWidth = max(width-prefixWidth, 1)
	}

	var segments []string
	if text.Wrap && wrapWidth > 0 {
		segments = wrapText(text.Value, wrapWidth)
	} else {
		segments = []string{text.Value}
	}

	lines := make([]string, len(segments))
	for i, segment := range segments {
		linePrefix := prefix
		if i > 0 {
			linePrefix = continuation
		}

		content := linePrefix + segment
		if text.Truncate && width > 0 {
			content = truncateString(content, width, text.TruncateSuffix)
		}

		line := style.Render(content)
		if width > 0 {
			line = lipgloss.PlaceHorizontal(width, mapHorizontal(text.Align), line)
		}

		lines[i] = line
	}

	return strings.Join(lines, "\n")
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

func truncateString(text string, width int, suffix string) string {
	if width <= 0 || lipgloss.Width(text) <= width {
		return text
	}

	if suffix == "" {
		suffix = "..."
	}

	suffixWidth := lipgloss.Width(suffix)
	if suffixWidth >= width {
		return truncateRunes(text, width)
	}

	target := width - suffixWidth
	var builder strings.Builder
	currentWidth := 0

	for _, r := range text {
		rw := lipgloss.Width(string(r))
		if currentWidth+rw > target {
			break
		}
		builder.WriteRune(r)
		currentWidth += rw
	}

	return builder.String() + suffix
}

func truncateRunes(text string, width int) string {
	if width <= 0 {
		return ""
	}

	var builder strings.Builder
	currentWidth := 0

	for _, r := range text {
		rw := lipgloss.Width(string(r))
		if currentWidth+rw > width {
			break
		}
		builder.WriteRune(r)
		currentWidth += rw
	}

	return builder.String()
}

func joinWithSpacing(segments []string, spacing int, join func(lipgloss.Position, ...string) string, pos lipgloss.Position) string {
	if len(segments) == 0 {
		return ""
	}
	if spacing <= 0 || len(segments) == 1 {
		return join(pos, segments...)
	}

	spacer := lipgloss.NewStyle().Width(spacing).Render("")
	withSpacing := make([]string, 0, len(segments)*2-1)
	for i, segment := range segments {
		if i > 0 {
			withSpacing = append(withSpacing, spacer)
		}
		withSpacing = append(withSpacing, segment)
	}
	return join(pos, withSpacing...)
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
