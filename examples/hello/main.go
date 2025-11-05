package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sprucelabsai-community/bubbleviews"
	"github.com/sprucelabsai-community/bubbleviews/render"
)

type model struct {
	view  bubbleviews.View
	ready bool
}

func newModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.view = buildHelloView(msg.Width, msg.Height)
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "loading..."
	}

	return render.Render(m.view)
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func buildHelloView(width, height int) bubbleviews.View {
	return bubbleviews.View{
		Size: bubbleviews.Size{Width: width, Height: height},
		Children: []bubbleviews.Node{
			bubbleviews.BoxNode{
				Style: bubbleviews.BoxStyle{
					Border:      bubbleviews.BorderThick,
					BorderColor: bubbleviews.Color("63"),
					Padding: bubbleviews.Padding{
						Top:    1,
						Bottom: 1,
						Left:   2,
						Right:  2,
					},
					FillWidth:  true,
					FillHeight: true,
					HAlign:     bubbleviews.AlignCenter,
					VAlign:     bubbleviews.AlignCenter,
				},
				Content: bubbleviews.View{
					Children: []bubbleviews.Node{
						bubbleviews.BoxNode{
							Style: bubbleviews.BoxStyle{
								Border:      bubbleviews.BorderThin,
								BorderColor: bubbleviews.Color("205"),
								Padding: bubbleviews.Padding{
									Left:  2,
									Right: 2,
								},
							},
							Content: bubbleviews.View{
								Children: []bubbleviews.Node{
									bubbleviews.TextNode{
										Value: "Hello World",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
