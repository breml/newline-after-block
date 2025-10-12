package comments

import "fmt"

// Test cases for block statements followed by comments

func ifStatementWithCommentNoNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("next statement")
}

func ifStatementWithCommentAndNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}

	// This comment has proper spacing
	fmt.Println("next statement")
}

func ifElseWithCommentNoNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("next statement")
}

func ifElseWithCommentAndNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	}

	// This comment has proper spacing
	fmt.Println("next statement")
}

func forLoopWithCommentNoNewline() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("after loop")
}

func forLoopWithCommentAndNewline() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// This comment has proper spacing
	fmt.Println("after loop")
}

func rangeLoopWithCommentNoNewline() {
	items := []int{1, 2, 3}
	for _, item := range items {
		fmt.Println(item)
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("after loop")
}

func rangeLoopWithCommentAndNewline() {
	items := []int{1, 2, 3}
	for _, item := range items {
		fmt.Println(item)
	}

	// This comment has proper spacing
	fmt.Println("after loop")
}

func switchWithCommentNoNewline() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("after switch")
}

func switchWithCommentAndNewline() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	}

	// This comment has proper spacing
	fmt.Println("after switch")
}

func typeSwitchWithCommentNoNewline() {
	a := any("hello")
	switch v := a.(type) {
	case string:
		fmt.Println("string:", v)
	case int:
		fmt.Println("int:", v)
	default:
		fmt.Println("unknown type")
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("after type switch")
}

func typeSwitchWithCommentAndNewline() {
	a := any("hello")
	switch v := a.(type) {
	case string:
		fmt.Println("string:", v)
	case int:
		fmt.Println("int:", v)
	default:
		fmt.Println("unknown type")
	}

	// This comment has proper spacing
	fmt.Println("after type switch")
}

func selectWithCommentNoNewline() {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println(v)
	default:
		fmt.Println("default")
	} // want "missing newline after block statement"
	// This comment should have a blank line above
	fmt.Println("after select")
}

func selectWithCommentAndNewline() {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println(v)
	default:
		fmt.Println("default")
	}

	// This comment has proper spacing
	fmt.Println("after select")
}

func multiLineCommentNoNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	/*
		This is a multi-line comment
		that should have a blank line above
	*/
	fmt.Println("next statement")
}

func multiLineCommentWithNewline() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}

	/*
		This is a multi-line comment
		with proper spacing
	*/
	fmt.Println("next statement")
}

func blockWithCommentOnly() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	// Comment without following statement is still a violation
}

func blockWithCommentAndNewlineOnly() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	}

	// Comment with proper spacing, even without following statement
}

// Test exclusions: composite literals should NOT be flagged

func compositeLiteralWithComment() {
	person := struct {
		name string
		age  int
	}{
		name: "John",
		age:  30,
	}
	// This comment after a composite literal should NOT be flagged
	fmt.Println(person)
}

func sliceLiteralWithComment() {
	items := []int{
		1,
		2,
		3,
	}
	// This comment after a slice literal should NOT be flagged
	fmt.Println(items)
}

func mapLiteralWithComment() {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	// This comment after a map literal should NOT be flagged
	fmt.Println(m)
}

func inlineTypeDefinition() {
	type Company struct {
		Name string
	}
	// This comment after a struct type definition should NOT be flagged.
}

func nestedBlocksWithComments() {
	x := 5
	if x > 0 {
		if x > 3 {
			fmt.Println("greater than 3")
		} // want "missing newline after block statement"
		// Inner comment needs spacing
		fmt.Println("positive")
	} // want "missing newline after block statement"
	// Outer comment needs spacing
	fmt.Println("done")
}

func nestedBlocksWithCommentsCorrect() {
	x := 5
	if x > 0 {
		if x > 3 {
			fmt.Println("greater than 3")
		}

		// Inner comment with proper spacing
		fmt.Println("positive")
	}

	// Outer comment with proper spacing
	fmt.Println("done")
}

func blockFollowedByStatementThenComment() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	fmt.Println("next statement")
	// This comment is after a statement, not directly after a block - OK
}

func blockWithInlineCommentThenBlockComment() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	// Block comment - should have newline above
	fmt.Println("next")
}

func blockWithInlineCommentThenBlockCommentCorrect() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // inline comment

	// Block comment with proper spacing
	fmt.Println("next")
}
