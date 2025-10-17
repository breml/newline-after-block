package newlineafterblock_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	newlineafterblock "github.com/breml/newline-after-block"
)

func TestAnalyzer(t *testing.T) {
	analyzer := newlineafterblock.New()

	// Exclude the _excluded.go file for this test
	err := analyzer.Flags.Set("exclude", `.*_excluded\.go`)
	if err != nil {
		t.Fatalf("failed to set exclude flag: %v", err)
	}

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "blockstatements")
}

func TestAnalyzerStructLiterals(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "structliterals")
}

func TestAnalyzerComments(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "comments")
}

func TestAnalyzerCaseClauses(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "caseclauses")
}

func TestAnalyzerWithFixes(t *testing.T) {
	analyzer := newlineafterblock.New()

	// Exclude the _excluded.go file for this test
	err := analyzer.Flags.Set("exclude", `.*_excluded\.go`)
	if err != nil {
		t.Fatalf("failed to set exclude flag: %v", err)
	}

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "blockstatements")
}

func TestAnalyzerStructLiteralsWithFixes(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "structliterals")
}

func TestAnalyzerCommentsWithFixes(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "comments")
}

func TestAnalyzerCaseClausesWithFixes(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "caseclauses")
}

func TestAnalyzerAnonymousFuncs(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "anonymousfuncs")
}

func TestAnalyzerAnonymousFuncsWithFixes(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "anonymousfuncs")
}

func TestAnalyzerDeferPattern(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "deferpattern")
}

func TestAnalyzerDeferPatternWithFixes(t *testing.T) {
	analyzer := newlineafterblock.New()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "deferpattern")
}
