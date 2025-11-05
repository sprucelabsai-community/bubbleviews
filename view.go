package bubbleviews

import "strings"

// View describes a rectangular region containing zero or more children.
type View struct {
	Size     Size
	Children []Node
}

// Node represents a renderable element in the view tree.
type Node interface {
	isNode()
}

// Size represents a width and height measured in terminal cells.
type Size struct {
	Width  int
	Height int
}

// BoxNode draws a bordered container that can host another view.
type BoxNode struct {
	Style   BoxStyle
	Content View
}

func (BoxNode) isNode() {}

// BoxStyle captures border, padding, fill, and alignment rules for a box.
type BoxStyle struct {
	Border      BorderStyle
	BorderColor Color
	Padding     Padding
	FillWidth   bool
	FillHeight  bool
	HAlign      Alignment
	VAlign      Alignment
}

// FlexNode arranges child nodes along a single axis.
type FlexNode struct {
	Direction FlexDirection
	Spacing   int
	Items     []FlexItem
}

func (FlexNode) isNode() {}

// FlexItem references a node within a Flex layout.
type FlexItem struct {
	Node   Node
	Width  int // used when Direction == FlexDirectionRow
	Height int // used when Direction == FlexDirectionColumn
}

// FlexDirection expresses whether a FlexNode lays out children in a row or column.
type FlexDirection int

const (
	FlexDirectionRow FlexDirection = iota
	FlexDirectionColumn
)

// TextNode renders raw text with optional formatting.
type TextNode struct {
	Value              string
	Color              Color
	Bold               bool
	Wrap               bool
	Align              Alignment
	Prefix             string
	ContinuationPrefix string
}

func (TextNode) isNode() {}

// Padding expresses the inset around content.
type Padding struct {
	Top, Right, Bottom, Left int
}

// BorderStyle enumerates supported frame borders.
type BorderStyle string

const (
	BorderNone  BorderStyle = "none"
	BorderThin  BorderStyle = "normal"
	BorderThick BorderStyle = "thick"
)

// Alignment describes horizontal or vertical placement.
type Alignment string

const (
	AlignStart  Alignment = "start"
	AlignCenter Alignment = "center"
	AlignEnd    Alignment = "end"
)

// Color is a free-form string keyed by the renderer.
type Color string

// ListView generates a vertical flex list of bullet items.
type ListView struct {
	Title      string
	TitleColor Color
	ItemColor  Color
	Bullet     string
	Items      []string
	Spacing    int
}

// Node returns the flex node representing the list content.
func (l ListView) Node() Node {
	items := make([]FlexItem, 0, len(l.Items)+1)

	if strings.TrimSpace(l.Title) != "" {
		items = append(items, FlexItem{
			Node: TextNode{
				Value: l.Title,
				Color: l.TitleColor,
				Bold:  true,
			},
		})
	}

	bullet := l.Bullet
	if bullet == "" {
		bullet = "- "
	}
	continuation := strings.Repeat(" ", len([]rune(bullet)))

	for _, item := range l.Items {
		items = append(items, FlexItem{
			Node: TextNode{
				Value:              item,
				Color:              l.ItemColor,
				Wrap:               true,
				Prefix:             bullet,
				ContinuationPrefix: continuation,
			},
		})
	}

	return FlexNode{
		Direction: FlexDirectionColumn,
		Spacing:   l.Spacing,
		Items:     items,
	}
}
