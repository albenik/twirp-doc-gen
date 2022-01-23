package markdown_test

import (
	"testing"

	md "github.com/albenik-go/twirp-doc-gen/internal/markdown"
)

func TestUL_Markdown(t *testing.T) {
	runTestCases(t, []*testCase{{
		Name: "Single",
		Block: md.UL(
			md.T("Item"),
		),
		Result: "* Item\n",
	}, {
		Name: "Multi",
		Block: md.UL(
			md.T("Item 1"),
			md.T("Item 2"),
			md.T("Item 3"),
		),
		Result: "* Item 1\n* Item 2\n* Item 3\n",
	}, {
		Name: "Links",
		Block: md.UL(
			md.Link("#anchor1", "Link1"),
			md.Link("#anchor2", "Link2"),
			md.Link("#anchor3", "Link3"),
		),
		Result: "* [Link1](#anchor1)\n* [Link2](#anchor2)\n* [Link3](#anchor3)\n",
	}})
}

func TestOL_Markdown(t *testing.T) {
	runTestCases(t, []*testCase{{
		Name: "Single",
		Block: md.OL(
			md.T("Item"),
		),
		Result: "1. Item\n",
	}, {
		Name: "Multi",
		Block: md.OL(
			md.T("Item 1"),
			md.T("Item 2"),
			md.T("Item 3"),
		),
		Result: "1. Item 1\n2. Item 2\n3. Item 3\n",
	}})
}
