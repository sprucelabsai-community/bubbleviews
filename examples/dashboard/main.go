package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sprucelabsai-community/bubbleviews"
	"github.com/sprucelabsai-community/bubbleviews/render"
)

type cameraStatus struct {
	id          string
	name        string
	fps         int
	dropped     int
	lastMessage string
}

type statusState struct {
	booting   bool
	cameras   []cameraStatus
	selected  selection
	lastEvent string
}

type selection struct {
	inSummary bool
	cameraIdx int
}

type model struct {
	width  int
	height int
	state  statusState
}

func newModel() *model {
	return &model{
		state: statusState{
			booting:   true,
			cameras:   []cameraStatus{},
			selected:  selection{inSummary: true, cameraIdx: -1},
			lastEvent: "Booting recorder service…",
		},
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

type addCameraMsg struct{}

type removeCameraMsg struct{}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.advanceFocus(1)
		case "shift+tab":
			m.advanceFocus(-1)
		case "enter":
			if m.state.selected.inSummary {
				return m, func() tea.Msg { return addCameraMsg{} }
			}
			if idx := m.state.selected.cameraIdx; idx >= 0 && idx < len(m.state.cameras) {
				return m, func() tea.Msg { return removeCameraMsg{} }
			}
		}
	case tickMsg:
		m.updateMetrics()
		return m, tea.Tick(time.Second, func(time.Time) tea.Msg { return tickMsg{} })
	case addCameraMsg:
		m.addCamera()
	case removeCameraMsg:
		m.removeFocusedCamera()
	}

	return m, nil
}

func (m *model) advanceFocus(direction int) {
	if m.state.selected.inSummary {
		if len(m.state.cameras) == 0 {
			if direction > 0 {
				m.addCamera()
			}
			return
		}
		if direction > 0 {
			m.state.selected.inSummary = false
			m.state.selected.cameraIdx = 0
		}
		return
	}

	if len(m.state.cameras) == 0 {
		m.state.selected.inSummary = true
		m.state.selected.cameraIdx = -1
		return
	}

	idx := m.state.selected.cameraIdx + direction
	if idx < 0 {
		m.state.selected.inSummary = true
		m.state.selected.cameraIdx = -1
		return
	}
	if idx >= len(m.state.cameras) {
		m.state.selected.inSummary = true
		m.state.selected.cameraIdx = -1
		return
	}
	m.state.selected.cameraIdx = idx
}

func (m *model) addCamera() {
	id := fmt.Sprintf("cam-%d", time.Now().UnixNano())
	cam := cameraStatus{
		id:          id,
		name:        fmt.Sprintf("Camera %d", len(m.state.cameras)+1),
		fps:         30,
		dropped:     0,
		lastMessage: "Receiving stream",
	}
	m.state.cameras = append(m.state.cameras, cam)
	m.state.selected.inSummary = false
	m.state.selected.cameraIdx = len(m.state.cameras) - 1
	m.state.booting = false
	m.state.lastEvent = fmt.Sprintf("Started %s", cam.name)
}

func (m *model) removeFocusedCamera() {
	idx := m.state.selected.cameraIdx
	if idx < 0 || idx >= len(m.state.cameras) {
		return
	}
	removed := m.state.cameras[idx]
	m.state.cameras = append(m.state.cameras[:idx], m.state.cameras[idx+1:]...)
	m.state.lastEvent = fmt.Sprintf("Stopped %s", removed.name)
	if len(m.state.cameras) == 0 {
		m.state.selected.inSummary = true
		m.state.selected.cameraIdx = -1
		return
	}
	if idx >= len(m.state.cameras) {
		m.state.selected.cameraIdx = len(m.state.cameras) - 1
	}
}

func (m *model) updateMetrics() {
	if len(m.state.cameras) == 0 {
		return
	}
	idx := rand.Intn(len(m.state.cameras))
	cam := m.state.cameras[idx]
	cam.fps = 25 + rand.Intn(11)
	if rand.Float64() < 0.2 {
		cam.dropped += rand.Intn(3)
		cam.lastMessage = "Detected minor packet loss"
	} else {
		cam.lastMessage = "Streaming nominal"
	}
	m.state.cameras[idx] = cam
}

func (m *model) View() string {
	var items []bubbleviews.FlexItem

	if len(m.state.cameras) > 0 {
		summary := buildSummaryRow(m.state)
		items = append(items, bubbleviews.FlexItem{Node: summary})
	}

	grid := buildCameraGrid(m.state)
	items = append(items, bubbleviews.FlexItem{Node: grid})

	layout := bubbleviews.FlexNode{
		Direction: bubbleviews.FlexDirectionColumn,
		Spacing:   1,
		Items:     items,
	}

	view := bubbleviews.View{
		Size:     bubbleviews.Size{Width: m.width, Height: m.height},
		Children: []bubbleviews.Node{layout},
	}

	return render.Render(view)
}

func buildSummaryRow(state statusState) bubbleviews.Node {
	statusLine := "Bridge online"
	if state.booting {
		statusLine = "Booting recorder service…"
	}
	cameraCount := fmt.Sprintf("Active cameras: %d", len(state.cameras))

	message := state.lastEvent
	if message == "" {
		message = "Waiting for events…"
	}

	flexItems := []bubbleviews.FlexItem{
		{
			Node: bubbleviews.TextNode{
				Value: "Streaming Recorder Service",
				Bold:  true,
			},
		},
		{
			Node: bubbleviews.TextNode{
				Value: statusLine,
				Color: bubbleviews.Color("36"),
			},
		},
		{
			Node: bubbleviews.TextNode{
				Value: cameraCount,
				Color: bubbleviews.Color("63"),
			},
		},
		{
			Node: bubbleviews.TextNode{
				Value: message,
				Color: bubbleviews.Color("244"),
			},
		},
	}

	addButton := bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThin,
			BorderColor: buttonBorderColor(state.selected.inSummary),
			Padding:     bubbleviews.Padding{Left: 2, Right: 2},
		},
		Content: bubbleviews.View{
			Children: []bubbleviews.Node{
				bubbleviews.TextNode{
					Value: "+ Add Camera",
					Bold:  state.selected.inSummary,
				},
			},
		},
	}

	flexItems = append(flexItems, bubbleviews.FlexItem{Node: addButton})

	return bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThick,
			BorderColor: bubbleviews.Color("63"),
			Padding:     bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2},
			FillWidth:   true,
		},
		Content: bubbleviews.View{
			Children: []bubbleviews.Node{
				bubbleviews.FlexNode{
					Direction: bubbleviews.FlexDirectionColumn,
					Spacing:   1,
					Items:     flexItems,
				},
			},
		},
	}
}

func buildCameraGrid(state statusState) bubbleviews.Node {
	if len(state.cameras) == 0 {
		placeholder := bubbleviews.BoxNode{
			Style: bubbleviews.BoxStyle{
				Border:      bubbleviews.BorderThin,
				BorderColor: bubbleviews.Color("240"),
				Padding:     bubbleviews.Padding{Top: 2, Bottom: 2, Left: 4, Right: 4},
				FillWidth:   true,
				FillHeight:  true,
				HAlign:      bubbleviews.AlignCenter,
				VAlign:      bubbleviews.AlignCenter,
			},
			Content: bubbleviews.View{
				Children: []bubbleviews.Node{
					bubbleviews.TextNode{
						Value: "No active cameras detected.",
						Color: bubbleviews.Color("244"),
					},
					bubbleviews.TextNode{
						Value: "Press TAB or Enter to add a camera.",
						Color: bubbleviews.Color("62"),
					},
					bubbleviews.BoxNode{
						Style: bubbleviews.BoxStyle{
							Border:      bubbleviews.BorderThin,
							BorderColor: buttonBorderColor(true),
							Padding:     bubbleviews.Padding{Left: 4, Right: 4},
						},
						Content: bubbleviews.View{
							Children: []bubbleviews.Node{
								bubbleviews.TextNode{
									Value: "+ Add Camera",
									Bold:  true,
								},
							},
						},
					},
				},
			},
		}

		return bubbleviews.BoxNode{
			Style: bubbleviews.BoxStyle{
				Border:      bubbleviews.BorderThin,
				BorderColor: bubbleviews.Color("63"),
				Padding:     bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2},
				FillWidth:   true,
				FillHeight:  true,
			},
			Content: bubbleviews.View{Children: []bubbleviews.Node{placeholder}},
		}
	}

	items := make([]bubbleviews.Node, 0, len(state.cameras))
	for idx, cam := range state.cameras {
		items = append(items, buildCameraPanel(cam, idx == state.selected.cameraIdx && !state.selected.inSummary))
	}

	gridContents := bubbleviews.FlowNode{
		ItemMinWidth: 28,
		ItemSpacing:  4,
		RowSpacing:   1,
		Items:        items,
	}

	return bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThin,
			BorderColor: bubbleviews.Color("63"),
			Padding:     bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2},
			FillWidth:   true,
			FillHeight:  true,
		},
		Content: bubbleviews.View{
			Children: []bubbleviews.Node{gridContents},
		},
	}
}

func buildCameraPanel(cam cameraStatus, focused bool) bubbleviews.Node {
	header := bubbleviews.TextNode{
		Value: cam.name,
		Bold:  true,
	}

	metrics := bubbleviews.FlexNode{
		Direction: bubbleviews.FlexDirectionColumn,
		Spacing:   0,
		Items: []bubbleviews.FlexItem{
			{Node: bubbleviews.TextNode{Value: fmt.Sprintf("FPS: %d", cam.fps)}},
			{Node: bubbleviews.TextNode{Value: fmt.Sprintf("Dropped frames: %d", cam.dropped)}},
			{Node: bubbleviews.TextNode{Value: cam.lastMessage, Color: bubbleviews.Color("244")}},
		},
	}

	removeButton := bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThin,
			BorderColor: buttonBorderColor(focused),
			Padding:     bubbleviews.Padding{Left: 2, Right: 2},
		},
		Content: bubbleviews.View{
			Children: []bubbleviews.Node{
				bubbleviews.TextNode{
					Value: "Remove",
					Bold:  focused,
				},
			},
		},
	}

	return bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThin,
			BorderColor: bubbleviews.Color("63"),
			Padding:     bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2},
		},
		Content: bubbleviews.View{
			Children: []bubbleviews.Node{
				header,
				metrics,
				removeButton,
			},
		},
	}
}

func buttonBorderColor(focused bool) bubbleviews.Color {
	if focused {
		return bubbleviews.Color("205")
	}
	return bubbleviews.Color("240")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
