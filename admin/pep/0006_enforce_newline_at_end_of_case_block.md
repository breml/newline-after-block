# PEP 0006: Enforce Newline at End of Case Block

Case blocks in Go (in `switch` and `select` statements) should be followed by a newline to enhance code readability.
There is no blank line required after the last case block if it is followed by the closing brace of the switch/select statement.

## Examples

### Incorrect

```go
switch value {
case 1:
    doSomething() // missing blank line at end of case block
case 2:
    doSomethingElse() // missing blank line at end of case block
default:
    doDefault()
}
```

### Correct

```go
switch value {
case 1:
    doSomething() // Blank line present

case 2:
    doSomethingElse() // Blank line present

default:
    doDefault()
}
```
