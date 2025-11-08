# PEP 0008: Report Error at End of Anonymous Functions

Currently, the `newline-after-block` linter does not report missing new lines after anonymous functions.

## Example

### Current Behavior

No error is reported for the following code:

```go
package comments

import "fmt"

func wsl() {
    _ = func() {
        _ = 1
    }
    // Comment group comment[0]

    // Comment group comment[1]
    fmt.Println("")
}
```

### Desired Behavior

```go
package comments

import "fmt"

func wsl() {
    _ = func() {
        _ = 1
    } // want "missing newline after block statement"
    // Comment group comment[0]

    // Comment group comment[1]
    fmt.Println("")
}
```

Golden file:

```go
package comments

import "fmt"

func wsl() {
    _ = func() {
        _ = 1
    } // want "missing newline after block statement"

    // Comment group comment[0]

    // Comment group comment[1]
    fmt.Println("")
}
```
