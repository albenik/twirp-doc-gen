package markdown_test

import (
	"testing"

	md "github.com/albenik/twirp-doc-gen/internal/markdown"
)

func TestHeader_Markdown(t *testing.T) {
	t.Parallel()

	runTestCases(t, []*testCase{{
		Name:   "H1",
		Block:  md.TH1("Header"),
		Result: "# Header\n",
	}, {
		Name:   "H2",
		Block:  md.TH2("Header"),
		Result: "## Header\n",
	}, {
		Name:   "H3",
		Block:  md.TH3("Header"),
		Result: "### Header\n",
	}, {
		Name:   "H4",
		Block:  md.TH4("Header"),
		Result: "#### Header\n",
	}, {
		Name:   "H5",
		Block:  md.TH5("Header"),
		Result: "##### Header\n",
	}, {
		Name:   "H6",
		Block:  md.TH6("Header"),
		Result: "###### Header\n",
	}, {
		Name:   "H1_WithCode",
		Block:  md.H1(md.T("H1 "), md.Code("const Result = \"OK\"")),
		Result: "# H1 `const Result = \"OK\"`\n",
	}})
}
