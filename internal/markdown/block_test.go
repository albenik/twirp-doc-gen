package markdown_test

import (
	"testing"

	md "github.com/albenik/twirp-doc-gen/internal/markdown"
)

func TestGroup_Markdown(t *testing.T) {
	t.Parallel()

	runTestCases(t, []*testCase{{
		Name: "InlineBlocks",
		Block: md.G(
			md.T("Foo "),
			md.TI("Bar"),
			md.T(" "),
			md.TB("Baz"),
			md.T(" "),
			md.Code("code"),
		),
		Result: "Foo *Bar* **Baz** `code`",
	}, {
		Name: "NewlineBlocks",
		Block: md.G(
			md.TH1("H1"),
			md.TH2("H2"),
			md.TH3("H3"),
			md.TH4("H4"),
			md.TH5("H5"),
			md.TH6("H6"),
		),
		Result: "# H1\n\n## H2\n\n### H3\n\n#### H4\n\n##### H5\n\n###### H6\n",
	}, {
		Name: "SingleBlockInGroup",
		Block: md.G(
			md.TH1("H1"),
		),
		Result: "# H1\n",
	}})
}
