package markdown

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	NoAlign = iota
	AlignLeft
	AlignCenter
	AlignRight
)

var (
	tableColSep = []byte("|")
	hyphen      = []byte("-")
	colon       = []byte(":")
	space       = []byte(" ")
	br          = []byte("<br/>")
)

// CellP renders a paragraph for table cell.
func CellP(blocks ...Block) Block {
	return Wrap(nil, br, blocks...)
}

type Table struct {
	columns []*tableHeader
	rows    [][]Block
}

type tableHeader struct {
	Title []byte
	Align uint8
}

func (t *Table) AddColumn(title string, align uint8) {
	t.columns = append(t.columns, &tableHeader{
		Title: []byte(title),
		Align: align,
	})
}

func (t *Table) AppendRow(cols ...Block) {
	t.rows = append(t.rows, cols)
}

func (t *Table) AppendRowT(cols ...string) {
	blocks := make([]Block, 0, len(cols))
	for _, s := range cols {
		blocks = append(blocks, T(s))
	}
	t.AppendRow(blocks...)
}

func (t *Table) Markdown(w io.Writer) error { //nolint:funlen,gocognit,gocyclo,cyclop
	const (
		minWidth   = 3
		emptyWidth = 5
	)

	headers := make([]*tableHeader, len(t.columns))
	copy(headers, t.columns)

	widths := make([]int, len(t.columns))
	for i, col := range headers {
		colWidth := utf8.RuneCount(col.Title)
		if colWidth < minWidth {
			colWidth = minWidth
		}
		widths[i] = colWidth
	}

	rows := make([][][]byte, 0, len(t.rows))
	for _, row := range t.rows {
		cols := make([][]byte, 0, len(row))
		for i, col := range row {
			if i > len(headers) {
				headers = append(headers, &tableHeader{Title: space}) //nolint:makezero
				widths = append(widths, emptyWidth)                   //nolint:makezero
			}

			buf := bytes.NewBuffer(nil)
			if err := col.Markdown(buf); err != nil {
				return err
			}
			cols = append(cols, buf.Bytes())

			if wd := utf8.RuneCount(buf.Bytes()); wd > widths[i] {
				widths[i] = wd
			}
		}
		rows = append(rows, cols)
	}

	if _, err := w.Write(tableColSep); err != nil {
		return err
	}

	for i, hdr := range headers {
		align := padLeft

		switch hdr.Align {
		case AlignCenter:
			align = center
		case AlignRight:
			align = right
		}

		if _, err := w.Write(space); err != nil {
			return err
		}
		if _, err := w.Write(align(hdr.Title, widths[i])); err != nil {
			return err
		}
		if _, err := w.Write(space); err != nil {
			return err
		}
		if _, err := w.Write(tableColSep); err != nil {
			return err
		}
	}

	if _, err := w.Write(newline); err != nil {
		return err
	}

	if _, err := w.Write(tableColSep); err != nil {
		return err
	}

	for i, col := range headers {
		switch col.Align {
		case AlignLeft:
			if _, err := w.Write(colon); err != nil {
				return err
			}
			if _, err := w.Write(bytes.Repeat(hyphen, widths[i]+1)); err != nil {
				return err
			}
		case AlignCenter:
			if _, err := w.Write(colon); err != nil {
				return err
			}
			if _, err := w.Write(bytes.Repeat(hyphen, widths[i])); err != nil {
				return err
			}
			if _, err := w.Write(colon); err != nil {
				return err
			}
		case AlignRight:
			if _, err := w.Write(bytes.Repeat(hyphen, widths[i]+1)); err != nil {
				return err
			}
			if _, err := w.Write(colon); err != nil {
				return err
			}
		default:
			if _, err := w.Write(bytes.Repeat(hyphen, widths[i]+2)); err != nil { //nolint:gomnd
				return err
			}
		}

		if _, err := w.Write(tableColSep); err != nil {
			return err
		}
	}

	if _, err := w.Write(newline); err != nil {
		return err
	}

	for _, row := range rows {
		if _, err := w.Write(tableColSep); err != nil {
			return err
		}

		for i, col := range row {
			align := padLeft

			switch headers[i].Align {
			case AlignCenter:
				align = center
			case AlignRight:
				align = right
			}

			if _, err := w.Write(space); err != nil {
				return err
			}
			if _, err := w.Write(align(col, widths[i])); err != nil {
				return err
			}
			if _, err := w.Write(space); err != nil {
				return err
			}
			if _, err := w.Write(tableColSep); err != nil {
				return err
			}
		}

		if _, err := w.Write(newline); err != nil {
			return err
		}
	}

	return nil
}

func padLeft(b []byte, w int) []byte {
	return []byte(fmt.Sprintf("%*s", -w, string(b)))
}

func center(b []byte, w int) []byte {
	s := string(b)
	return []byte(fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))) //nolint:gomnd
}

func right(b []byte, w int) []byte {
	return []byte(fmt.Sprintf("%*s", w, string(b)))
}
