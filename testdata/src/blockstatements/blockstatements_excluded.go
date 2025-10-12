package blockstatements

import "fmt"

// This file contains block statement violations that are excluded from analysis
// using the --exclude flag in tests. The violations here demonstrate that
// excluded files are not checked by the linter.

func ifStatementViolation() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("next statement")
}

func forLoopViolation() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("after loop")
}

func rangeLoopViolation() {
	items := []int{1, 2, 3}
	for _, item := range items {
		fmt.Println(item)
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("after loop")
}

func switchStatementViolation() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	default:
		fmt.Println("other")
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("after switch")
}

func selectStatementViolation() {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println(v)
	default:
		fmt.Println("default")
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("after select")
}

func ifElseViolation() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("not positive")
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("next statement")
}

func nestedViolations() {
	x := 5
	y := 10
	if x > 0 {
		if y > 0 {
			fmt.Println("both positive")
		}
		fmt.Println("x positive")
	} // no "missing newline after block statement" report, since the file is excluded.
	fmt.Println("after outer if")
}
