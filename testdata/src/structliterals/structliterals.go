package structliterals

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Address struct {
	Street string
	City   string
}

func structLiteralWithoutNewline() {
	// Struct literals should not require newlines
	p := Person{
		Name: "John",
		Age:  30,
	}
	fmt.Println(p)
}

func structLiteralMultipleWithoutNewline() {
	p := Person{
		Name: "John",
		Age:  30,
	}
	a := Address{
		Street: "Main St",
		City:   "NYC",
	}
	fmt.Println(p, a)
}

func arrayLiteralWithoutNewline() {
	// Array/slice literals should not require newlines
	arr := []int{
		1,
		2,
		3,
	}
	fmt.Println(arr)
}

func mapLiteralWithoutNewline() {
	// Map literals should not require newlines
	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	fmt.Println(m)
}

func nestedStructLiteral() {
	type Company struct {
		Name    string
		Address Address
	}

	c := Company{
		Name: "ACME",
		Address: Address{
			Street: "Main St",
			City:   "NYC",
		},
	}
	fmt.Println(c)
}

func mixedBlockAndLiteral() {
	x := 5
	if x > 0 { // want "missing newline after block statement"
		fmt.Println("positive")
	}
	p := Person{
		Name: "John",
		Age:  30,
	}
	fmt.Println(p)
}
