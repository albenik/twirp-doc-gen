package markdown_test

import (
	"testing"

	md "github.com/albenik/twirp-doc-gen/internal/markdown"
)

func TestRenderTable(t *testing.T) {
	t.Parallel()

	runTestCases(t, []*testCase{{
		Name: "SingleColumnNoAlign",
		Block: func() *md.Table {
			table := new(md.Table)
			table.AddColumn("H", md.NoAlign)
			table.AppendRowT("A")
			table.AppendRowT("B")
			table.AppendRowT("C")
			return table
		}(),
		Result: "| H   |\n|-----|\n| A   |\n| B   |\n| C   |\n",
	}, {
		Name: "ThreeColumns",
		Block: func() *md.Table {
			table := new(md.Table)
			table.AddColumn("H", md.NoAlign)
			table.AddColumn("H", md.AlignLeft)
			table.AddColumn("H", md.AlignCenter)
			table.AddColumn("H", md.AlignRight)
			table.AppendRowT("1", "1", "1", "1")
			table.AppendRowT("22", "22", "22", "22")
			table.AppendRowT("333", "333", "333", "333")
			table.AppendRowT("44444", "44444", "44444", "44444")
			return table
		}(),
		Result: "| H     | H     |   H   |     H |\n" +
			"|-------|:------|:-----:|------:|\n" +
			"| 1     | 1     |   1   |     1 |\n" +
			"| 22    | 22    |  22   |    22 |\n" +
			"| 333   | 333   |  333  |   333 |\n" +
			"| 44444 | 44444 | 44444 | 44444 |\n",
	}, {
		Name: "DifferentWidths",
		Block: func() *md.Table {
			table := new(md.Table)
			table.AddColumn("1", md.NoAlign)
			table.AddColumn("12345", md.AlignLeft)
			table.AddColumn("1234567", md.AlignCenter)
			table.AddColumn("1", md.AlignRight)
			table.AppendRowT("1234567", "123", "1", "1")
			table.AppendRowT("1", "12345", "12", "1")
			table.AppendRowT("1", "12", "123", "1234567")
			return table
		}(),
		Result: "| 1       | 12345 | 1234567 |       1 |\n" +
			"|---------|:------|:-------:|--------:|\n" +
			"| 1234567 | 123   |    1    |       1 |\n" +
			"| 1       | 12345 |   12    |       1 |\n" +
			"| 1       | 12    |   123   | 1234567 |\n",
	}, {
		Name: "Formatting",
		Block: func() *md.Table {
			table := new(md.Table)
			table.AddColumn("A", md.AlignLeft)
			table.AddColumn("B", md.AlignCenter)
			table.AddColumn("C", md.AlignRight)
			table.AppendRow(md.Code("code"), md.Code("code"), md.Code("code"))
			table.AppendRow(md.T("12345678"), md.T("12345678"), md.T("12345678"))
			table.AppendRow(md.TB("bold"), md.TI("italic"), md.T("plain"))
			return table
		}(),
		Result: "| A        |    B     |        C |\n" +
			"|:---------|:--------:|---------:|\n" +
			"| `code`   |  `code`  |   `code` |\n" +
			"| 12345678 | 12345678 | 12345678 |\n" +
			"| **bold** | *italic* |    plain |\n",
	}})
}
