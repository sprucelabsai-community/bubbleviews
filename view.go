package bubbleviews

// View describes a rectangular region containing zero or more children.
type View struct {
	Size     Size
	Children []ViewChild
}

// ViewChild wraps the supported concrete components.
type ViewChild struct {
	Frame     *FrameView
	Placement *PlacementView
	Button    *ButtonView
	Columns   *ColumnsView
	List      *ListView
}

// Size represents a width and height measured in terminal cells.
type Size struct {
	Width  int
	Height int
}

// FrameView draws a bordered container that can host another view.
type FrameView struct {
	Border      BorderStyle
	BorderColor Color
	Padding     Padding
	FillWidth   bool
	FillHeight  bool
	Content     *View
}

// PlacementView positions child content within a bounded area.
type PlacementView struct {
	AreaWidth  int
	AreaHeight int
	HAlign     Alignment
	VAlign     Alignment
	Content    *View
}

// ColumnsView arranges child views horizontally.
type ColumnsView struct {
	Spacing int
	Columns []ColumnView
}

// ColumnView represents a single column within a ColumnsView.
type ColumnView struct {
	Width   int // absolute width; 0 means share remaining space evenly
	Content *View
}

// ButtonView renders a clickable or highlightable label.
type ButtonView struct {
	Label       string
	Border      BorderStyle
	BorderColor Color
	Padding     Padding
	TextColor   Color
}

// Padding expresses the inset around content.
type Padding struct {
	Top, Right, Bottom, Left int
}

// ListView renders a titled bulleted list.
type ListView struct {
	Title      string
	TitleColor Color
	ItemColor  Color
	Bullet     string
	Items      []string
}

// BorderStyle enumerates supported frame and button borders.
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
