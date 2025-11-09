# Even Rows Example

- **Scenario:** Platform status board where cards lock to three per row and reserve explicit width percentages.
- **Primary struct:** `bubbleviews.EvenRowGrid` built in `buildGridView`, which emits `EqualWidthRow`-style flex groups automatically.

```go
grid := bubbleviews.EvenRowGrid{
    Items: []bubbleviews.EvenRowItem{
        {Node: buildStatusCard(edgeRecorder), WidthPercent: 50},
        {Node: buildStatusCard(cloudRelay), WidthPercent: 25},
        {Node: buildStatusCard(reviewStations), WidthPercent: 25},
        // ...
    },
    MaxPerRow:  3,
    Spacing:    2,
    RowSpacing: 1,
}.Node()
```

### What this tests
- Percentage-based widths via `FlexItem.WidthPercent` along with automatic remainder handling.
- Uniform row heights (`RowHeight`) plus per-card `FillWidth` boxes that stretch to the given allocation.
- Text truncation (ellipsis) whenever content exceeds its slot, illustrating the `Truncate` flag.

### Run it
```sh
go run ./examples/even_rows
```
