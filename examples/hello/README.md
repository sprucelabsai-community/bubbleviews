# Hello Example

- **Scenario:** Minimal full-screen layout that fills the terminal, centers a child box, and exercises nested `BoxNode` styling.
- **Primary struct:** `bubbleviews.View` built by `buildHelloView`, featuring a thick outer `BoxNode` containing a thin bordered child.

```go
func buildHelloView(width, height int) bubbleviews.View {
    return bubbleviews.View{
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
}
```

### What this tests
- `FillWidth`/`FillHeight` handling when the parent view knows its size.
- Alignment logic inside `BoxNode` (`HAlign` + `VAlign` centering).
- Basic text rendering without wrapping or truncation.

### Run it
```sh
go run ./examples/hello
```
