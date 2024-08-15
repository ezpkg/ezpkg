package one

// +genz:sample 10

// this is comment of A
//
// +genz:a this directive should be ignored from comment text
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

// +genz:b this directive should be ignored
type B int

// +genz:last 20: number:int * x
