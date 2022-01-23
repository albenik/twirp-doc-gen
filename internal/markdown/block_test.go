package markdown_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	md "github.com/albenik-go/twirp-doc-gen/internal/markdown"
)

func TestGroup_Markdown(t *testing.T) {
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

type testCase struct {
	Name   string
	Block  md.Block
	Result string
}

func runTestCases(t *testing.T, cases []*testCase) {
	t.Helper()

	for _, c := range cases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			require.NoError(t, c.Block.Markdown(buf))
			require.Equal(t, c.Result, buf.String())
		})
	}
}
