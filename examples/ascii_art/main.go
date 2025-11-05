package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sprucelabsai-community/bubbleviews"
	"github.com/sprucelabsai-community/bubbleviews/render"
)

type model struct {
	width, height int
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
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	lines := []string{
		`  ____   _ _           _        _         _`,
		` / __ \ (_) |         | |      (_)       | |`,
		`| |  | |_ _| |__  _ __| |_ _ __ _  ____ _| |___`,
		`| |  | | | | '_ \| '__| __| '__| |/ / _' | / __|`,
		`| |__| | | | |_) | |  | |_| |  | /\\ (_| | \__ \\`,
		` \____/|_|_|_.__/|_|   \__|_|  |_/_\\__,_|_|___/`,
		``,
		`   Odin Analytics`,
	}

	art := bubbleviews.ASCIIArtNode{
		Lines: lines,
		Align: bubbleviews.AlignCenter,
		Bold:  true,
		Color: bubbleviews.Color("63"),
	}
	view := bubbleviews.View{
		Size: bubbleviews.Size{Width: m.width, Height: m.height},
		Children: []bubbleviews.Node{
			bubbleviews.BoxNode{
				Style: bubbleviews.BoxStyle{
					Border:      bubbleviews.BorderThick,
					BorderColor: bubbleviews.Color("63"),
					Padding:     bubbleviews.Padding{Top: 2, Bottom: 2, Left: 4, Right: 4},
					FillWidth:   true,
					FillHeight:  true,
					HAlign:      bubbleviews.AlignCenter,
					VAlign:      bubbleviews.AlignCenter,
				},
				Content: bubbleviews.View{Children: []bubbleviews.Node{art}},
			},
		},
	}

	return render.Render(view)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
