# PEP 0009: Allow Defer After Block Statements

It is a common pattern in Go to execute a command and immediately defer its cleanup.

For example:

```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
```

Similarly, when working with files, it is common to open a file and immediately defer its closure:

```go
file, err := os.Open("example.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()
```

However, the current implementation of `newline-after-block` enforces a newline
after block statements, which can lead to less idiomatic code in these scenarios:

```go
file, err := os.Open("example.txt")
if err != nil {
    log.Fatal(err)
}

defer file.Close()
```

This proposal suggests to allow `defer` statements immediately following block
statements without requiring a newline in between if, and only if, the block
statement is checking for an error. This change would enhance code readability
and maintain idiomatic Go practices.

## Specification

### Scope

This exception applies to `if` statements only (not `for`, `switch`, or `select`).

### Pattern Detection

The analyzer detects error-checking patterns using type-based detection:

- The condition is a binary expression with the `!=` operator
- One operand is a variable whose type implements the `error` interface
- The other operand is `nil`

This approach allows the analyzer to recognize error checks regardless of the variable name,
supporting both standard patterns like `if err != nil` and variations like `if e != nil`,
`if problem != nil`, or any other variable name, as long as the type implements `error`.

### Behavior

1. **Single defer after error check**: No blank line required between the `if` block and the `defer` statement

   ```go
   file, err := os.Open("example.txt")
   if err != nil {
       return err
   }
   defer file.Close()

   processFile(file)  // blank line required here
   ```

2. **Multiple consecutive defer statements**: No blank line required between consecutive `defer` statements

   ```go
   file, err := os.Open("example.txt")
   if err != nil {
       return err
   }
   defer file.Close()
   defer mu.Unlock()

   processFile(file)  // blank line required after all defers
   ```

3. **Defer followed by non-defer statement**: A blank line IS required after the last `defer` before any non-`defer` statement

   ```go
   file, err := os.Open("example.txt")
   if err != nil {
       return err
   }
   defer file.Close()

   data, err := readFile(file)  // blank line required
   ```

4. **Defer followed by another block statement**: A blank line IS required before the next block statement

   ```go
   file, err := os.Open("example.txt")
   if err != nil {
       return err
   }
   defer file.Close()

   if someOtherCondition {  // blank line required
       doSomething()
   }
   ```

The `defer` statement(s) are considered part of the same "paragraph" as the
preparation statement and error-checking block, forming a cohesive unit.
