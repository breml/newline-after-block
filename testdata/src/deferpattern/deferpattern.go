package deferpattern

import (
	"fmt"
	"os"
)

// Test 1: Basic pattern - if err != nil followed by defer (should NOT warn)
func basicDeferAfterErrorCheck() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("processing file")
	return nil
}

// Test 2: Multiple consecutive defers after error check (should NOT warn)
func multipleConsecutiveDefers() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	defer fmt.Println("cleanup")

	fmt.Println("processing file")
	return nil
}

// Test 3: Defer followed by regular statement without blank line (SHOULD warn)
func deferFollowedByStatementNoBlankLine() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close() // want "missing newline after block statement"
	data := []byte("test")

	fmt.Println(data)
	return nil
}

// Test 4: Defer followed by regular statement with blank line (should NOT warn)
func deferFollowedByStatementWithBlankLine() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	data := []byte("test")

	fmt.Println(data)
	return nil
}

// Test 5: Defer followed by block statement without blank line (SHOULD warn)
func deferFollowedByBlockNoBlankLine() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close() // want "missing newline after block statement"
	if true {
		fmt.Println("in block")
	}

	return nil
}

// Test 6: Defer followed by block statement with blank line (should NOT warn)
func deferFollowedByBlockWithBlankLine() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	if true {
		fmt.Println("in block")
	}

	return nil
}

// Test 7: Non-error-check if followed by defer (SHOULD warn)
func nonErrorCheckIfFollowedByDefer() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	defer fmt.Println("cleanup")

	fmt.Println("done")
}

// Test 8: if err == nil (wrong operator) followed by defer (SHOULD warn)
func wrongOperatorFollowedByDefer() error {
	file, err := os.Open("example.txt")
	if err == nil {
		fmt.Println("success")
	} // want "missing newline after block statement"
	defer fmt.Println("cleanup")

	fmt.Println(file)
	return nil
}

// Test 9: if with different variable name followed by defer (should NOT warn - type-based detection)
func differentVariableNameFollowedByDefer() error {
	file, e := os.Open("example.txt")
	if e != nil {
		return e
	}
	defer file.Close()

	fmt.Println(file)
	return nil
}

// Test 10: Nested - defer after error check inside another block (should NOT warn)
func nestedDeferAfterErrorCheck() {
	if true {
		file, err := os.Open("example.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Println("processing")
	}
}

// Test 11: Multiple defers followed by statement without blank line (SHOULD warn)
func multipleDefersThenStatementNoBlankLine() error {
	file, err := os.Open("example.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	defer fmt.Println("cleanup") // want "missing newline after block statement"
	data := []byte("test")

	fmt.Println(data)
	return nil
}

// Test 12: if nil != err pattern (reversed operands) followed by defer (should NOT warn)
func reversedOperandsFollowedByDefer() error {
	file, err := os.Open("example.txt")
	if nil != err {
		return err
	}
	defer file.Close()

	fmt.Println("processing file")
	return nil
}

// Test 13: Standalone defer statements (should NOT warn between them)
func standaloneConsecutiveDefers() {
	defer fmt.Println("first")
	defer fmt.Println("second")
	defer fmt.Println("third")

	fmt.Println("body")
}

// Test 14: Standalone defer followed by statement without blank line (SHOULD warn)
func standaloneDeferFollowedByStatementNoBlankLine() {
	defer fmt.Println("cleanup") // want "missing newline after block statement"
	x := 5

	fmt.Println(x)
}

// Test 15: Standalone defer followed by statement with blank line (should NOT warn)
func standaloneDeferFollowedByStatementWithBlankLine() {
	defer fmt.Println("cleanup")

	x := 5

	fmt.Println(x)
}

// Test 16: Custom error type with different name followed by defer (should NOT warn - type-based detection)
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func customErrorTypeFollowedByDefer() error {
	var myErr error = &customError{msg: "test error"}
	if myErr != nil {
		return myErr
	}
	defer fmt.Println("cleanup")

	fmt.Println("success")
	return nil
}

// Test 17: Error variable with unusual name followed by defer (should NOT warn - type-based detection)
func unusualErrorNameFollowedByDefer() error {
	file, problem := os.Open("example.txt")
	if problem != nil {
		return problem
	}
	defer file.Close()

	fmt.Println("processing file")
	return nil
}

// Test 18: Non-error type with != nil check followed by defer (SHOULD warn)
func nonErrorTypeFollowedByDefer() {
	var ptr *int
	if ptr != nil {
		fmt.Println("not nil")
	} // want "missing newline after block statement"
	defer fmt.Println("cleanup")

	fmt.Println("done")
}

// Test 19: Reversed operands with different variable name (should NOT warn - type-based detection)
func reversedOperandsWithDifferentName() error {
	file, problem := os.Open("example.txt")
	if nil != problem {
		return problem
	}
	defer file.Close()

	fmt.Println("processing file")
	return nil
}
