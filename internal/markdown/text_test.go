package markdown_test

import (
	"testing"

	md "github.com/albenik-go/twirp-doc-gen/internal/markdown"
)

func TestText_Markdown(t *testing.T) {
	runTestCases(t, []*testCase{{
		Name:   "Plain",
		Block:  md.T("Test OK"),
		Result: "Test OK",
	}, {
		Name:   "Italic",
		Block:  md.I(md.T("Test OK")),
		Result: "*Test OK*",
	}, {
		Name:   "Bold",
		Block:  md.B(md.T("Test OK")),
		Result: "**Test OK**",
	}, {
		Name:   "CombinedFormatting",
		Block:  md.B(md.TI("Test"), md.T(" OK")),
		Result: "***Test* OK**",
	}, {
		Name:   "Paragraph",
		Block:  md.TP("Test OK"),
		Result: "Test OK\n",
	}, {
		Name: "MultiParagraph",
		Block: md.G(
			md.TP("P#1"),
			md.TP("P#2"),
			md.TP("P#3"),
		),
		Result: "P#1\n\nP#2\n\nP#3\n",
	}})
}
