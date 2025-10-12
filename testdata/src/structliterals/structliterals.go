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

func inlineTypeDefinition() {
	type Company struct {
		Name    string
		Address Address
	}
	fmt.Println(Company{})
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

func structLiteralBeforeIf() {
	p := Person{
		Name: "John",
		Age:  30,
	}
	if p.Age > 18 { // want "missing newline after block statement"
		fmt.Println("adult")
	}
	fmt.Println(p.Name)
}

func arrayLiteralBeforeIf() {
	arr := []int{
		1,
		2,
		3,
	}
	if len(arr) > 0 { // want "missing newline after block statement"
		fmt.Println("not empty")
	}
	fmt.Println(arr)
}

func mapLiteralBeforeIf() {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	if len(m) > 0 { // want "missing newline after block statement"
		fmt.Println("not empty")
	}
	fmt.Println(m)
}

func sliceLiteralBeforeIf() {
	s := []string{
		"a",
		"b",
	}
	if len(s) > 0 { // want "missing newline after block statement"
		fmt.Println("not empty")
	}
	fmt.Println(s)
}
