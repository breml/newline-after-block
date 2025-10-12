// Command newline-after-block is a linter that checks for newlines after block statements.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	newlineafterblock "github.com/breml/newline-after-block"
)

func main() {
	singlechecker.Main(newlineafterblock.New())
}
