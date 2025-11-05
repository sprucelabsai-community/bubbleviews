package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sprucelabsai-community/bubbleviews"
	"github.com/sprucelabsai-community/bubbleviews/render"
)

type model struct {
	width        int
	height       int
	selected     int
	clickMessage string
}

var commands = []string{
	"Add camera column",
	"Remove focused camera",
	"Toggle debug overlay",
}

func newModel() model {
	return model{
		clickMessage: "Press enter to invoke a command.",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(commands)-1 {
				m.selected++
			}
		case "enter":
			m.clickMessage = fmt.Sprintf("Command executed: %s", commands[m.selected])
		}
	}

	return m, nil
}

func (m model) View() string {
	view := bubbleviews.View{
		Size: bubbleviews.Size{Width: m.width, Height: m.height},
		Children: []bubbleviews.Node{
			buildHeaderRow(),
			buildCommandList(m),
		},
	}

	return render.Render(view)
}

func buildHeaderRow() bubbleviews.Node {
	return bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThick,
			BorderColor: bubbleviews.Color("63"),
			FillWidth:   true,
		},
		Content: bubbleviews.View{
			Children: []bubbleviews.Node{
				bubbleviews.TextNode{
					Value: "Interactive Controls Demo",
					Bold:  true,
				},
				bubbleviews.TextNode{
					Value: "Use ↑/↓ to change selection. Press enter to fire the focused command. q to quit.",
					Color: bubbleviews.Color("244"),
				},
			},
		},
	}
}

func buildCommandList(m model) bubbleviews.Node {
	items := make([]bubbleviews.FlexItem, 0, len(commands)+1)

	for i, cmd := range commands {
		isSelected := i == m.selected
		borderColor := bubbleviews.Color("240")
		if isSelected {
			borderColor = bubbleviews.Color("205")
		}

		item := bubbleviews.FlexItem{
			Node: bubbleviews.BoxNode{
				Style: bubbleviews.BoxStyle{
					Border:      bubbleviews.BorderThin,
					BorderColor: borderColor,
					Padding: bubbleviews.Padding{
						Left:   2,
						Right:  2,
						Top:    0,
						Bottom: 0,
					},
					FillWidth: true,
				},
				Content: bubbleviews.View{
					Children: []bubbleviews.Node{
						bubbleviews.TextNode{
							Value: cmd,
							Bold:  isSelected,
						},
					},
				},
			},
		}

		items = append(items, item)
	}

	items = append(items, bubbleviews.FlexItem{
		Node: bubbleviews.TextNode{
			Value: m.clickMessage,
			Color: bubbleviews.Color("36"),
		},
	})

	return bubbleviews.BoxNode{
		Style: bubbleviews.BoxStyle{
			Border:      bubbleviews.BorderThin,
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
				bubbleviews.FlexNode{
					Direction: bubbleviews.FlexDirectionColumn,
					Spacing:   1,
					Items:     items,
				},
			},
		},
	}
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
