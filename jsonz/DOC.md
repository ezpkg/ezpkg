# ezpkg.io/jsonz

Package [jsonz](https://pkg.go.dev/ezpkg.io/jsonz) is a minimal implementation of json parser and transformer in Go. The
`Parse()` function returns an iter over the JSON object, which can be used to traverse the JSON object.

## Examples

### 1. Iterate over the JSON object:

```go
package main

import (
	"fmt"
	"ezpkg.io/jsonz"
)

func main() {
	jsonStr := `{"name": "Alice", "age": 24, "address": {"city": "The Sun", "zip": 10101}}`

	fmt.Println("| Path | Index | Key | Token |")
	fmt.Println("|------|-------|-----|-------|")
	for item, err := range jsonz.Parse([]byte(jsonStr)) {
		if err != nil {
			panic(err)
		}
		fmt.Printf("| %v | %v | %v | %v |\n", item.GetPath(), item.Index, item.Key, item.Token)
	}
}
```

Will output:

| Path         | Index | Key       | Token     |
|--------------|-------|-----------|-----------|
|              | 0     |           | {         |
| name         | 0     | "name"    | "Alice"   |
| age          | 1     | "age"     | 24        |
| scores       | 2     | "scores"  | [         |
| scores.0     | 0     |           | 9         |
| scores.1     | 1     |           | 10        |
| scores.2     | 2     |           | 8         |
| scores       | 2     |           | ]         |
| address      | 3     | "address" | {         |
| address.city | 0     | "city"    | "The Sun" |
| address.zip  | 1     | "zip"     | 10101     |
| address      | 3     |           | }         |
|              | 0     |           | }         |

### 2. Reconstruction of the JSON object:

You can reconstruct the JSON object by iterating over the JSON object and adding commas between the tokens. This is useful when you want to modify the JSON object and write it back to a file.

```go
package main

import (
	"bytes"
	"fmt"
	"ezpkg.io/jsonz"
)

func main() {
	jsonStr := `{"name": "Alice", "age": 24, "address": {"city": "The Sun", "zip": 10101}}`
	
	var b bytes.Buffer
	var lastTokenType jsonz.TokenType
	for item, err := range jsonz.Parse([]byte(jsonStr)) {
		if err != nil {
			panic(err)
		}
		if jsonz.ShouldAddComma(lastTokenType, item.Token.Type()) {
			b.WriteByte(',')
		}
		if item.Key.IsValue() {
			b.Write(item.Key.Raw())
			b.WriteByte(':')
		}
		b.Write(item.Token.Raw())
		lastTokenType = item.Token.Type()
	}
	fmt.Printf("%s\n", b.Bytes())
}
```

Will output:

```json
{"name":"Alice","age":24,"scores":[9,10,8],"address":{"city":"The Sun","zip":10101}}
```

### 3. Reformat the JSON object:

This example will reformat the JSON object by adding newlines and indentation.

```go
package main

import (
	"fmt"
	"ezpkg.io/jsonz"
	"ezpkg.io/bytez"
)

func main() {
	jsonStr := `{"name": "Alice", "age": 24, "address": {"city": "The Sun", "zip": 10101}}`
	
	var b bytez.Buffer
	var lastTokenType jsonz.TokenType
	for item, err := range jsonz.Parse([]byte(jsonStr)) {
		if err != nil {
			panic(err)
		}
		if jsonz.ShouldAddComma(lastTokenType, item.Token.Type()) {
			b.Print(",")
		}
		b.Println()
		for i := 0; i < item.Level; i++ {
			b.Print("  ")
		}
		if item.Key.IsValue() {
			b.WriteZ(item.Key.Raw())
			b.Print(": ")
		}
		b.WriteZ(item.Token.Raw())
		lastTokenType = item.Token.Type()
	}
	fmt.Printf("%s\n", b.Bytes())
}
```

Will output:

```json
{
  "name": "Alice",
  "age": 24,
  "scores": [
	9,
	10,
	8
  ],
  "address": {
	"city": "The Sun",
	"zip": 10101
  }
}
```
