package newlineafterblock_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	newlineafterblock "github.com/breml/newline-after-block"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, newlineafterblock.Analyzer, "blockstatements")
}

func TestAnalyzerStructLiterals(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, newlineafterblock.Analyzer, "structliterals")
}
