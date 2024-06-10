Packages testingz provides utilities for testing. Support ignoring spaces and using placeholders.

## Examples

```go
formatted, isDiff := testingz.DiffByChar("code: 123A", "code: ███A")
// isDiff: true

formatted, isDiff := testingz.DiffByCharZ("code: 123A", "code: ███A")
// isDiff: false

// placeholder is useful for comparing tests with uuid or random values
formatted, isDiff := testingz.DiffByLineZ(left, right)
left := "id: ████\ncode: AF███\nname: Alice\n"
right := "id: 1234\ncode: AF123\nname: Alice\n"
// isDiff: false
```
