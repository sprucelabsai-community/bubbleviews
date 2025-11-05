# bubbleviews Agent Brief

## Guiding Principles
- Treat Bubble Tea as a thin host for event handling only. All layout logic lives in a pure render model defined in `view.go`.
- The render model must remain “dumb”: immutable data structs with no behaviour or framework imports. Renderers consume the model and produce output.
- Rendering happens in one shot. Build the full `bubbleviews.View`, hand it to a renderer (currently Lip Gloss), and return the rendered string.
- Examples drive development. We vibe-code new features by adding `examples/<name>` programs instead of formal unit tests.

## Render Model Expectations
- `View` nodes can be nested via `ViewChild` unions. Each child populates exactly one concrete component (`Frame`, `Placement`, `Button`, etc.).
- Components describe intent only: borders use symbolic enums, colours are arbitrary strings, alignment is semantic (`AlignStart`, `AlignCenter`, `AlignEnd`).
- Renderers are responsible for translating the model into terminal output. Keep renderer logic pure and side-effect free where possible.

## Project Workflow
- Add new demos under `examples/` to showcase behaviour and validate the render pipeline. Use `go run ./examples/<name>` (with `GOCACHE=$(pwd)/.gocache` if needed).
- Skip unit tests unless explicitly requested; document any noteworthy quirks or manual test steps alongside the examples.
- For debugging or live preview, use the VS Code launch config “Run: Hello World Example” which opens an integrated terminal.

## Collaboration Notes
- Prefer expanding the render model over sprinkling conditional logic in renderers.
- When introducing new visual patterns, define fresh component structs (e.g. `ListView`, `TextBlockView`) and adapt the renderer accordingly.
- Keep generated output ASCII-only unless a specific feature warrants otherwise.
- Update this brief if the overall philosophy shifts, so future agents stay aligned.
