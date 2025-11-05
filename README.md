# bubbleviews

Pure data render models for Bubble Tea + Lip Gloss. Define your interface as Go
structs, hand it to a renderer, and let Bubble Tea do nothing more than ferry
events and repaint strings.

- **Render model first.** Compose Box, Flex, and Text nodes defined in `view.go`.
- **Renderer second.** `render.Render` walks the model and emits Lip Gloss markup.
- **Bubble Tea last.** The Bubble Tea program simply caches the latest render model.

<p align="center">
<a href="examples/hello"><code>examples/hello</code></a> ·
<a href="#render-model">Render Model</a> ·
<a href="#examples">Examples</a>
</p>

---

## Quick start

```sh
git clone git@github.com:sprucelabsai-community/bubbleviews.git
cd bubbleviews
GOCACHE=$(pwd)/.gocache go run ./examples/hello
```

Prefer VS Code? Use the built-in launch config **Run: Hello World Example** and
the demo will start in the integrated terminal (`q`, `esc`, or `ctrl+c` to exit).

---

## Render model

All layout intent lives in a set of plain structs:

```go
view := bubbleviews.View{
	Size: bubbleviews.Size{Width: width, Height: height},
	Children: []bubbleviews.Node{
		bubbleviews.BoxNode{
			Style: bubbleviews.BoxStyle{
				Border:      bubbleviews.BorderThick,
				BorderColor: bubbleviews.Color("63"),
				Padding:     bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2},
				FillWidth:   true,
				FillHeight:  true,
				HAlign:      bubbleviews.AlignCenter,
				VAlign:      bubbleviews.AlignCenter,
			},
			Content: bubbleviews.View{
				Children: []bubbleviews.Node{
					bubbleviews.BoxNode{
						Style: bubbleviews.BoxStyle{
							Border:      bubbleviews.BorderThin,
							BorderColor: bubbleviews.Color("205"),
							Padding:     bubbleviews.Padding{Left: 2, Right: 2},
						},
						Content: bubbleviews.View{
							Children: []bubbleviews.Node{
								bubbleviews.TextNode{Value: "Hello World"},
							},
						},
					},
				},
			},
		},
	},
}
```

Renderers remain pure, translating these intent structs into terminal output.
There are no Bubble Tea imports inside the render model, and the renderer never
mutates the model it receives.

---

## Examples

We’re vibe-coding this project—examples are our tests.

- [`examples/hello`](examples/hello): full-screen box with a centered “Hello World” button.
- [`examples/tasks`](examples/tasks): split-pane dashboard showing outstanding and completed checklists side-by-side.
- [`examples/interactivity`](examples/interactivity): keyboard-driven command list showcasing interactive focus/selection.


---

## License

MIT © Spruce Labs
