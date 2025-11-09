# ASCII Art Example

- **Scenario:** Centered ASCII banner with foreground color and bold weight.
- **Primary struct:** `bubbleviews.ASCIIArtNode` returned inside the top-level `bubbleviews.View`.

```go
art := bubbleviews.ASCIIArtNode{
    Lines: []string{
        "  ____   _ _           _        _         _",
        " / __ \\ (_) |         | |      (_)       | |",
        "| |  | |_ _| |__  _ __| |_ _ __ _  ____ _| |___",
        "| |  | | | | '_ \\| '__| __| '__| |/ / _' | / __|",
        "| |__| | | | |_) | |  | |_| |  | /\\ (_| | \\__ \\",
        " \\____/|_|_|_.__/|_|   \\__|_|  |_/_\\\\__,_|_|___/",
        "",
        "   Odin Analytics",
    },
    Align: bubbleviews.AlignCenter,
    Bold:  true,
    Color: bubbleviews.Color("63"),
}
```

### What this tests
- Rendering preformatted text blocks with alignment semantics.
- Applying Lip Gloss foreground colors and bold styling directly from the render model.
- Mix-and-match with plain `TextNode`s for captions under the art.

### Run it
```sh
go run ./examples/ascii_art
```
