package main

import (
	"fmt"

	"ezpkg.io/errorz"
	"ezpkg.io/jsonz"
)

type Address struct {
	HouseNumber int    `json:"house_number"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

func main() {
	{
		b := jsonz.NewBuilder("", "    ")
		// open an object
		b.Add("", jsonz.TokenObjectOpen)

		// add a few fields
		b.Add("name", "Alice")
		b.Add("age", 22)
		b.Add("email", "alice@example.com")
		b.Add("phone", "(+84) 123-456-789")

		// open an array
		b.Add("languages", jsonz.TokenArrayOpen)
		b.Add("", "English")
		b.Add("", "Vietnamese")
		b.Add("", jsonz.TokenArrayClose)
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
		b.Add("", jsonz.TokenObjectClose)

		out := errorz.Must(b.Bytes())
		fmt.Printf("\n--- build json ---\n%s\n", out)
	}
}
