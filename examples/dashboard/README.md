# Dashboard Example

- **Scenario:** Recorder status dashboard with a summary header and dynamic camera grid using `FlowNode`.
- **Primary struct:** `bubbleviews.FlowNode` returned by `buildCameraGrid`, which auto-wraps camera panels as terminal width changes.

```go
gridContents := bubbleviews.FlowNode{
    ItemMinWidth: 28,
    ItemSpacing:  4,
    RowSpacing:   1,
    Items:        panels,
}
```

### What this tests
- Large mixed-content views that combine ASCII art, text, and nested boxes.
- Responsive camera panels managed by `FlowNode` with `ItemMinWidth` so columns wrap as space changes.
- Integration with Bubble Tea state updates (add/remove cameras, periodic metric refresh).

### Run it
```sh
go run ./examples/dashboard
```
