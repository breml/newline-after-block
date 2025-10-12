# PEP 0005: Newline After Block Before Comment

If a block statement (if, for, switch, select, etc.) is followed by a comment,
there should be a blank line between the block and the comment.

This does not apply to comments, that follow composite literals or struct type definitions.

## Examples

### Incorrect

```go
if condition {
    // do something
}
// This is a comment about the next statement.
nextStatement()
```

### Correct

```go
if condition {
    // do something
}

// This is a comment about the next statement.
nextStatement()
```
