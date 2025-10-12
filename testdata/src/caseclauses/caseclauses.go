package caseclauses

import "fmt"

// Switch statements - violations
func switchWithoutNewlineBetweenCases() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one") // want "missing newline after case block"
	case 2:
		fmt.Println("two") // want "missing newline after case block"
	default:
		fmt.Println("other") // No blank line needed before }
	}
}

// Switch statements - correct
func switchWithNewlineBetweenCases() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")

	case 2:
		fmt.Println("two")

	default:
		fmt.Println("other") // No blank line needed before }
	}
}

// Single case - no violation
func switchSingleCase() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")
	}
}

// Empty case - should be skipped
func switchEmptyCase() {
	x := 1
	switch x {
	case 1:
	case 2:
		fmt.Println("two")
	}
}

// Multiple statements in case
func switchMultipleStatementsInCase() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")
		fmt.Println("still one") // want "missing newline after case block"
	case 2:
		fmt.Println("two")
	}
}

// Fallthrough cases
func switchFallthrough() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")
		fallthrough // want "missing newline after case block"
	case 2:
		fmt.Println("two")
	}
}

// Select statements - violations
func selectWithoutNewlineBetweenCases() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	select {
	case v := <-ch1:
		fmt.Println(v) // want "missing newline after case block"
	case ch2 <- 42:
		fmt.Println("sent") // want "missing newline after case block"
	default:
		fmt.Println("default")
	}
}

// Select statements - correct
func selectWithNewlineBetweenCases() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	select {
	case v := <-ch1:
		fmt.Println(v)

	case ch2 <- 42:
		fmt.Println("sent")

	default:
		fmt.Println("default")
	}
}

// Type switch - violations
func typeSwitchWithoutNewlineBetweenCases() {
	var x interface{} = "hello"
	switch v := x.(type) {
	case string:
		fmt.Println("string:", v) // want "missing newline after case block"
	case int:
		fmt.Println("int:", v) // want "missing newline after case block"
	default:
		fmt.Println("unknown") // No blank line needed before }
	}
}

// Type switch - correct
func typeSwitchWithNewlineBetweenCases() {
	var x interface{} = "hello"
	switch v := x.(type) {
	case string:
		fmt.Println("string:", v)

	case int:
		fmt.Println("int:", v)

	default:
		fmt.Println("unknown") // No blank line needed before }
	}
}

// Nested switches
func nestedSwitchWithoutNewlines() {
	x := 1
	y := 2
	switch x {
	case 1:
		switch y {
		case 1:
			fmt.Println("1,1") // want "missing newline after case block"
		case 2:
			fmt.Println("1,2")
		} // want "missing newline after case block"
	case 2:
		fmt.Println("x=2") // No blank line needed before }
	}
}

// Nested switches - correct
func nestedSwitchWithNewlines() {
	x := 1
	y := 2
	switch x {
	case 1:
		switch y {
		case 1:
			fmt.Println("1,1")

		case 2:
			fmt.Println("1,2") // No blank line needed before }
		}

	case 2:
		fmt.Println("x=2") // No blank line needed before }
	}
}

// Comments between cases - violation
func switchWithCommentNoNewline() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one") // want "missing newline after case block"
	// This comment needs a blank line above
	case 2:
		fmt.Println("two") // No blank line needed before }
	}
}

// Comments between cases - correct
func switchWithCommentAndNewline() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")

	// This comment has a blank line above
	case 2:
		fmt.Println("two") // No blank line needed before }
	}
}

// Multiple cases on same line (using comma)
func switchMultipleCasesOneLine() {
	x := 1
	switch x {
	case 1, 2, 3:
		fmt.Println("1, 2, or 3")

	case 4, 5:
		fmt.Println("4 or 5")

	default:
		fmt.Println("other") // No blank line needed before }
	}
}

// Switch with no default
func switchNoDefault() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")

	case 2:
		fmt.Println("two") // No blank line needed before }
	}
}

// Switch with expression
func switchWithExpression() {
	x := 5
	switch {
	case x < 0:
		fmt.Println("negative") // want "missing newline after case block"
	case x == 0:
		fmt.Println("zero") // want "missing newline after case block"
	case x > 0:
		fmt.Println("positive") // No blank line needed before }
	}
}

// Switch with initialization
func switchWithInit() {
	switch x := getValue(); x {
	case 1:
		fmt.Println("one")

	case 2:
		fmt.Println("two")

	default:
		fmt.Println("other") // No blank line needed before }
	}
}

func getValue() int {
	return 1
}

// Complex case with block statements inside
func switchWithBlocksInside() {
	x := 1
	switch x {
	case 1:
		if true {
			fmt.Println("one")
		}

		fmt.Println("still case 1") // want "missing newline after case block"
	case 2:
		for i := 0; i < 2; i++ {
			fmt.Println(i)
		}

		fmt.Println("case 2") // No blank line needed before }
	}
}
