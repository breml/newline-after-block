// Command newline-after-block is a linter that checks for newlines after block statements.
package main

import (
	newlineafterblock "github.com/breml/newline-after-block"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(newlineafterblock.New())
}
