# PEP 0002: Ensure Handling of Additional Edge Cases in newline-after-block Linter

This document outlines the additional edge cases that the `newline-after-block` linter must handle to ensure robust functionality.
The linter should correctly identify and enforce newline requirements after block statements while ignoring specific scenarios where
newlines are not necessary.

## Additional Edge Cases to Handle

### Newline required

- After `else` blocks
- Switch type statement

### No newline required

- `else` block at the end of a block.
- Newline after `type` definition.
- Struct, array, slice or map literal before `if` statement.

## Additional Test Cases

The following test cases should be added to the `testdata` directory to ensure the linter correctly handles these edge cases:

- Tests for all the above listed edge cases.
- Nested blocks like nested `if`, `for` or `switch` statements.
