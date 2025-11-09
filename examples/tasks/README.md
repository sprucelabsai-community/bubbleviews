# Tasks Example

- **Scenario:** Two-column task dashboard with headers, list bullets, and shared container borders.
- **Primary struct:** `bubbleviews.FlexNode` returned by `buildTasksView`, which arranges two boxed `ListView` nodes side-by-side.

```go
columns := bubbleviews.FlexNode{
    Direction: bubbleviews.FlexDirectionRow,
    Spacing:   4,
    Items: []bubbleviews.FlexItem{
        {Node: bubbleviews.BoxNode{Style: bubbleviews.BoxStyle{Border: bubbleviews.BorderThin, BorderColor: bubbleviews.Color("69"), Padding: bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2}, FillWidth: true, FillHeight: true}, Content: bubbleviews.View{Children: []bubbleviews.Node{outstandingList}}}},
        {Node: bubbleviews.BoxNode{Style: bubbleviews.BoxStyle{Border: bubbleviews.BorderThin, BorderColor: bubbleviews.Color("108"), Padding: bubbleviews.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2}, FillWidth: true, FillHeight: true}, Content: bubbleviews.View{Children: []bubbleviews.Node{completedList}}}},
    },
}
```

### What this tests
- Horizontal flex layouts with explicit spacing between columns.
- Reusable `ListView` helper emitting bullet-aligned `TextNode`s.
- Box padding + fill behavior inside a larger framed view.

### Run it
```sh
go run ./examples/tasks
```
