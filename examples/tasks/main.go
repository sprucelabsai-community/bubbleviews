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
		m.view = buildTasksView(msg.Width, msg.Height)
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

func buildTasksView(width, height int) bubbleviews.View {
	outstandingList := bubbleviews.ListView{
		Title:      "Outstanding Tasks",
		TitleColor: bubbleviews.Color("69"),
		Items: []string{
			"Review camera calibration",
			"Sync ingest pipeline with S3",
			"Schedule QA playback session",
			"Confirm alert latency thresholds",
		},
	}.Node()

	completedList := bubbleviews.ListView{
		Title:      "Completed Tasks",
		TitleColor: bubbleviews.Color("108"),
		ItemColor:  bubbleviews.Color("244"),
		Items: []string{
			"Deploy recorder v1.4.2",
			"Archive April footage",
			"Send daily digest to ops",
			"Reconcile storage billing",
		},
	}.Node()

	columns := bubbleviews.FlexNode{
		Direction: bubbleviews.FlexDirectionRow,
		Spacing:   4,
		Items: []bubbleviews.FlexItem{
			{
				Node: bubbleviews.BoxNode{
					Style: bubbleviews.BoxStyle{
						Border:      bubbleviews.BorderThin,
						BorderColor: bubbleviews.Color("69"),
						Padding: bubbleviews.Padding{
							Top:    1,
							Bottom: 1,
							Left:   2,
							Right:  2,
						},
						FillWidth:  true,
						FillHeight: true,
					},
					Content: bubbleviews.View{
						Children: []bubbleviews.Node{
							outstandingList,
						},
					},
				},
			},
			{
				Node: bubbleviews.BoxNode{
					Style: bubbleviews.BoxStyle{
						Border:      bubbleviews.BorderThin,
						BorderColor: bubbleviews.Color("108"),
						Padding: bubbleviews.Padding{
							Top:    1,
							Bottom: 1,
							Left:   2,
							Right:  2,
						},
						FillWidth:  true,
						FillHeight: true,
					},
					Content: bubbleviews.View{
						Children: []bubbleviews.Node{
							completedList,
						},
					},
				},
			},
		},
	}

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
				},
				Content: bubbleviews.View{
					Children: []bubbleviews.Node{
						columns,
					},
				},
			},
		},
	}
}
