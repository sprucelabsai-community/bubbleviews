# Equal Width Dashboard Example

- **Scenario:** Recorder dashboard remix that chunks camera cards into rows of three and lets each card claim identical width.
- **Primary struct:** `bubbleviews.EqualWidthRow` created inside `buildGridView`, with cards rendered by `buildStatusCard`.

```go
row := bubbleviews.EqualWidthRow{
    Items: []bubbleviews.Node{
        buildStatusCard(statuses[i]),
        buildStatusCard(statuses[i+1]),
        buildStatusCard(statuses[i+2]),
    },
    Spacing: 2,
}.Node()
```

### What this tests
- `FlexItem.Grow` weights to share parent width evenly without manual math.
- Combining `EqualWidthRow` rows inside a column flex to get responsive grids with deterministic layout.
- Showcasing truncation + ellipsis when long status copy overflows a card.

### Run it
```sh
go run ./examples/dashboard_equalrow
```
