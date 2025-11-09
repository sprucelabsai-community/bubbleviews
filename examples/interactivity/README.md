# Interactivity Example

- **Scenario:** Keyboard-driven command list that toggles focus and highlights the active action.
- **Primary struct:** `bubbleviews.FlexNode` assembled in `buildCommandList`, containing stacked `BoxNode` controls that reflect the model state.

```go
func buildCommandList(m model) bubbleviews.Node {
    items := make([]bubbleviews.FlexItem, 0, len(commands)+1)
    for i, cmd := range commands {
        isSelected := i == m.selected
        borderColor := bubbleviews.Color("240")
        if isSelected {
            borderColor = bubbleviews.Color("205")
        }
        items = append(items, bubbleviews.FlexItem{
            Node: bubbleviews.BoxNode{
                Style: bubbleviews.BoxStyle{Border: bubbleviews.BorderThin, BorderColor: borderColor, Padding: bubbleviews.Padding{Left: 2, Right: 2}, FillWidth: true},
                Content: bubbleviews.View{Children: []bubbleviews.Node{bubbleviews.TextNode{Value: cmd, Bold: isSelected}}},
            },
        })
    }
    items = append(items, bubbleviews.FlexItem{Node: bubbleviews.TextNode{Value: m.clickMessage, Color: bubbleviews.Color("36")}})
    return bubbleviews.BoxNode{Style: bubbleviews.BoxStyle{Border: bubbleviews.BorderThin, BorderColor: bubbleviews.Color("63"), Padding: bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2}, FillWidth: true, FillHeight: true}, Content: bubbleviews.View{Children: []bubbleviews.Node{bubbleviews.FlexNode{Direction: bubbleviews.FlexDirectionColumn, Spacing: 1, Items: items}}}}
}
```

### What this tests
- Bubble Tea events updating a cached render model without entangling UI code with the renderer.
- Conditional styling (focused vs. unfocused) expressed purely through `BoxStyle` and `TextNode` fields.
- Multi-line text wrapping with prefixes for command descriptions.

### Run it
```sh
go run ./examples/interactivity
```
