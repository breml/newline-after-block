# PEP 0003: Report Error at End of Block

Currently, the `newline-after-block` linter reports errors at the beginning of the block statement. This can be confusing for users
trying to identify the exact location of the issue. To improve usability, the linter should be modified to report errors at the end
of the block statement itself.

## Example

### Current Behavior

```go
if condition { // want "missing newline after block statement"
    doSomething()
}
nextStatement()
```

### Desired Behavior

```go
if condition {
    doSomething()
} // want "missing newline after block statement"
nextStatement()
```
