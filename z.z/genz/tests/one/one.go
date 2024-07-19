package one

// +ggen:sample 10

// this is comment of A
//
// +ggen:a this directive should be ignored from comment text
type A struct {
	Zero struct{}

	// comment of One
	One int

	Two string

	//
	// comment of Three
	//
	Three bool
}

// +ggen:b this directive should be ignored
type B int

// +ggen:last 20: number:int * x
