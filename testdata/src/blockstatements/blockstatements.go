package blockstatements

import "fmt"

func ifStatementWithoutNewline() {
	x := 5
	if x > 0 { // want "missing newline after block statement"
		fmt.Println("positive")
	}
	fmt.Println("next statement")
}

func ifStatementWithNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}

	fmt.Println("next statement")
}

func ifStatementAtEnd() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}
}

func ifElseStatement() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	}

	fmt.Println("next statement")
}

func ifElseIfStatement() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	}

	fmt.Println("next statement")
}

func forLoopWithoutNewline() {
	for i := 0; i < 5; i++ { // want "missing newline after block statement"
		fmt.Println(i)
	}
	fmt.Println("after loop")
}

func forLoopWithNewline() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	fmt.Println("after loop")
}

func forLoopAtEnd() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
}

func rangeLoopWithoutNewline() {
	items := []int{1, 2, 3}
	for _, item := range items { // want "missing newline after block statement"
		fmt.Println(item)
	}
	fmt.Println("after loop")
}

func rangeLoopWithNewline() {
	items := []int{1, 2, 3}
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("after loop")
}

func switchStatementWithoutNewline() {
	x := 2
	switch x { // want "missing newline after block statement"
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	}
	fmt.Println("after switch")
}

func switchStatementWithNewline() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	}

	fmt.Println("after switch")
}

func selectStatementWithoutNewline() {
	ch := make(chan int)
	select { // want "missing newline after block statement"
	case v := <-ch:
		fmt.Println(v)
	default:
		fmt.Println("default")
	}
	fmt.Println("after select")
}

func selectStatementWithNewline() {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println(v)
	default:
		fmt.Println("default")
	}

	fmt.Println("after select")
}

func multipleStatementsWithMixedViolations() {
	x := 5

	if x > 0 { // want "missing newline after block statement"
		fmt.Println("positive")
	}
	for i := 0; i < x; i++ {
		fmt.Println(i)
	}

	switch x {
	case 5:
		fmt.Println("five")
	}
}
