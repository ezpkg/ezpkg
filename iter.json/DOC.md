# ezpkg.io/iter.json

Package [iter.json](https://pkg.go.dev/ezpkg.io/iter.json) is JSON parser and transformer in Go. The `Parse()` function returns an iterator over the JSON object, which can be used to traverse the JSON object. And `Builder` can be used to construct a JSON object. Together, they provide a powerful way to iterate and manipulate your JSON data with Go [iterators](https://pkg.go.dev/iter#hdr-Iterators).

## Examples

Given an example [alice.json](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/alice.json) file:

```json
{
  "name": "Alice",
  "age": 24,
  "scores": [9, 10, 8],
  "address": {
    "city": "The Sun",
    "zip": 10101
  }
}
```

You can query and manipulate the JSON object in various ways:

### 1. Iterating JSON:

Use `for range Parse()` to iterate over a JSON data, then print the path, key, token, and level of each item. See [examples/01.iter](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/01.iter/main.go).

```go
package main

import (
    "fmt"

    "ezpkg.io/errorz"
    iterjson "ezpkg.io/iter.json"
)

func main() {
    data := `{"name": "Alice", "age": 24, "scores": [9, 10, 8], "address": {"city": "The Sun", "zip": 10101}}`

    // 🎄Example: iterate over json
    fmt.Printf("| %12v | %10v | %10v |%v|\n", "PATH", "KEY", "TOKEN", "LVL")
    fmt.Println("| ------------ | ---------- | ---------- | - |")
    for item, err := range iterjson.Parse([]byte(data)) {
        errorz.MustZ(err)

        fmt.Printf("| %12v | %10v | %10v | %v |\n", item.GetPathString(), item.Key, item.Token, item.Level)
    }
}
```

The code will output:

```text
|         PATH |        KEY |      TOKEN |LVL|
| ------------ | ---------- | ---------- | - |
|              |            |          { | 0 |
|         name |     "name" |    "Alice" | 1 |
|          age |      "age" |         24 | 1 |
|       scores |   "scores" |          [ | 1 |
|     scores.0 |            |          9 | 2 |
|     scores.1 |            |         10 | 2 |
|     scores.2 |            |          8 | 2 |
|       scores |            |          ] | 1 |
|      address |  "address" |          { | 1 |
| address.city |     "city" |  "The Sun" | 2 |
|  address.zip |      "zip" |      10101 | 2 |
|      address |            |          } | 1 |
|              |            |          } | 0 |
```

### 2. Building JSON:

Use `Builder` to build a JSON data. It accepts optional arguments for indentation. See [examples/02.builder](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/02.builder/main.go).

```go
b := iterjson.NewBuilder("", "    ")
// open an object
b.Add("", iterjson.TokenObjectOpen)

// add a few fields
b.Add("name", "Alice")
b.Add("age", 22)
b.Add("email", "alice@example.com")
b.Add("phone", "(+84) 123-456-789")

// open an array
b.Add("languages", iterjson.TokenArrayOpen)
b.Add("", "English")
b.Add("", "Vietnamese")
b.Add("", iterjson.TokenArrayClose)
// close the array

// accept any type that can marshal to json
b.Add("address", Address{
    HouseNumber: 42,
    Street:      "Ly Thuong Kiet",
    City:        "Ha Noi",
    Country:     "Vietnam",
})

// accept []byte as raw json
b.Add("pets", []byte(`[{"type":"cat","name":"Kitty","age":2},{"type":"dog","name":"Yummy","age":3}]`))

// close the object
b.Add("", iterjson.TokenObjectClose)

out := errorz.Must(b.Bytes())
fmt.Printf("\n--- build json ---\n%s\n", out)
```

Which will output the JSON with indentation:

```json
{
    "name": "Alice",
    "age": 22,
    "email": "alice@example.com",
    "phone": "(+84) 123-456-789",
    "languages": [
        "English",
        "Vietnamese"
    ],
    "address": {"house_number":42,"street":"Ly Thuong Kiet","city":"Ha Noi","country":"Vietnam"},
    "pets": [
        {
            "type": "cat",
            "name": "Kitty",
            "age": 2
        },
        {
            "type": "dog",
            "name": "Yummy",
            "age": 3
        }
    ]
}
```

### 3. Formatting JSON:

You can reconstruct or format a JSON data by sending its key and values to a `Builder`. See [examples/03.reformat](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/03.reformat/main.go).

```go
{
    // 🐝Example: minify json
    b := iterjson.NewBuilder("", "")
    for item, err := range iterjson.Parse(data) {
        errorz.MustZ(err)
        b.AddRaw(item.Key, item.Token)
    }
    out := errorz.Must(b.Bytes())
    fmt.Printf("\n--- minify ---\n%s\n----------\n", out)
}
{
    // 🦋Example: format json
    b := iterjson.NewBuilder("👉   ", "\t")
    for item, err := range iterjson.Parse(data) {
        errorz.MustZ(err)
        b.AddRaw(item.Key, item.Token)
    }
    out := errorz.Must(b.Bytes())
    fmt.Printf("\n--- reformat ---\n%s\n----------\n", out)
}
```

The first example minifies the JSON while the second example formats it with prefix "👉" on each line.

```text
--- minify ---
{"name":"Alice","age":24,"scores":[9,10,8],"address":{"city":"The Sun","zip":10101}}
----------

--- reformat ---
👉   {
👉       "name": "Alice",
👉       "age": 24,
👉       "scores": [
👉           9,
👉           10,
👉           8
👉       ],
👉       "address": {
👉           "city": "The Sun",
👉           "zip": 10101
👉       }
👉   }
----------
```

### 4. Adding line numbers

In this example, we add line numbers to the JSON output, by adding a `b.WriteNewline()` before the `fmt.Fprintf()` call. See [examples/04.line_number](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/04.line_number/main.go).

```go
// 🐞Example: print with line number
i := 0
b := iterjson.NewBuilder("", "    ")
for item, err := range iterjson.Parse(data) {
    i++
    errorz.MustZ(err)
    b.WriteNewline(item.Token.Type())

    // 👉 add line number
    fmt.Fprintf(b, "%3d    ", i)
    b.Add(item.Key, item.Token)
}
out := errorz.Must(b.Bytes())
fmt.Printf("\n--- line number ---\n%s\n----------\n", out)
```

This will output:

```text
  1    {
  2        "name": "Alice",
  3        "age": 24,
  4        "scores": [
  5            9,
  6            10,
  7            8
  8        ],
  9        "address": {
 10            "city": "The Sun",
 11            "zip": 10101
 12        }
 13    }
```

### 5. Adding comments

By putting a `fmt.Fprintf(comment)` between `b.WriteComma()` and `b.WriteNewline()`, you can add a comment to the end of each line. See [examples/05.comment](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/05.comment/main.go).

```go
i, newlineIdx, maxIdx := 0, 0, 30
b := iterjson.NewBuilder("", "    ")
for item, err := range iterjson.Parse(data) {
    errorz.MustZ(err)
    b.WriteComma(item.Token.Type())

    // 👉 add comment
    if i > 0 {
        length := b.Len() - newlineIdx
        fmt.Fprint(b, strings.Repeat(" ", maxIdx-length))
        fmt.Fprintf(b, "// %2d", i)
    }
    i++

    b.WriteNewline(item.Token.Type())
    newlineIdx = b.Len() // save the newline index

    b.Add(item.Key, item.Token)
}
length := b.Len() - newlineIdx
fmt.Fprint(b, strings.Repeat(" ", maxIdx-length))
fmt.Fprintf(b, "// %2d", i)

out := errorz.Must(b.Bytes())
fmt.Printf("\n--- comment ---\n%s\n----------\n", out)
```

This will output:

```text
{                             //  1
    "name": "Alice",          //  2
    "age": 24,                //  3
    "scores": [               //  4
        9,                    //  5
        10,                   //  6
        8                     //  7
    ],                        //  8
    "address": {              //  9
        "city": "The Sun",    // 10
        "zip": 10101          // 11
    }                         // 12
}                             // 13
```

### 6. Filtering JSON and extracting values

There are `item.GetPathString()` and `item.GetRawPath()` to get the path of the current item. You can use them to filter the JSON data. See [examples/06.filter_print](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/06.filter_print/main.go).

Example with `item.GetPathString()` and `regexp`:

```go
fmt.Printf("\n--- filter: GetPathString() ---\n")
i := 0
for item, err := range iterjson.Parse(data) {
    i++
    errorz.MustZ(err)

    path := item.GetPathString()
    switch {
    case path == "name",
        strings.Contains(path, "address"):
        // continue
    default:
        continue
    }

    // 👉 print with line number
    fmt.Printf("%2d %20s . %s\n", i, item.Token, item.GetPath())
}
```

Example with `item.GetRawPath()` and `path.Match()`:

```go
fmt.Printf("\n--- filter: GetRawPath() ---\n")
i := 0
for item, err := range iterjson.Parse(data) {
    i++
    errorz.MustZ(err)

    path := item.GetRawPath()
    switch {
    case path.Match("name"),
        path.Contains("address"):
        // continue
    default:
        continue
    }

    // 👉 print with line number
    fmt.Printf("%2d %20s . %s\n", i, item.Token, item.GetPath())
}
```

Both examples will output:

```text
 2              "Alice" . name
 9                    { . address
10            "The Sun" . address.city
11                10101 . address.zip
12                    } . address
```

### 7. Filtering JSON and returning a new JSON

By combining the `Builder` with the option `SetSkipEmptyStructures(false)` and the filtering logic, you can filter the JSON data and return a new JSON. See [examples/07.filter_json](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/07.filter_json/main.go)

```go
// 🦁Example: filter and output json
b := iterjson.NewBuilder("", "    ")
b.SetSkipEmptyStructures(true) // 👉 skip empty [] or {}
for item, err := range iterjson.Parse(data) {
    errorz.MustZ(err)
    if item.Token.IsOpen() || item.Token.IsClose() {
        b.Add(item.Key, item.Token)
        continue
    }

    path := item.GetPathString()
    switch {
    case path == "name",
        strings.Contains(path, "address"):
        // continue
    default:
        continue
    }

    b.Add(item.Key, item.Token)
}
out := errorz.Must(b.Bytes())
fmt.Printf("\n--- filter: output json ---\n%s\n----------\n", out)
```

This example will return a new JSON with only the filtered fields:

```json
{
    "name": "Alice",
    "address": {
        "city": "The Sun",
        "zip": 10101
    }
}
```

### 8. Editing values

This is an example for editing values in a JSON data. Assume that we are using number ids for our API. The ids are too big and JavaScript can't handle them. We need to convert them to strings. See [examples/08.number_id](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/08.number_id/main.go) and [order.json](https://github.com/ezpkg/ezpkg/blob/main/iter.json/examples/order.json).

Iterate over the JSON data, find all `_id` fields and convert the number ids to strings:

```go
b := iterjson.NewBuilder("", "    ")
for item, err := range iterjson.Parse(data) {
    errorz.MustZ(err)
    key, _ := item.GetRawPath().Last().ObjectKey()
    if strings.HasSuffix(key, "_id") {
        id, err0 := item.Token.GetInt()
        if err0 == nil {
            b.Add(item.Key, strconv.Itoa(id))
            continue
        }
    }
    b.Add(item.Key, item.Token)
}
out := errorz.Must(b.Bytes())
fmt.Printf("\n--- convert number id ---\n%s\n----------\n", out)
```

This will add quotes to the number ids:

```json
{
    "order_id": "12345678901234",
    "number": 12,
    "customer_id": "12345678905678",
    "items": [
        {
            "item_id": "12345678901042",
            "quantity": 1,
            "price": 123.45
        },
        {
            "item_id": "12345678901098",
            "quantity": 2,
            "price": 234.56
        }
    ]
}
```
