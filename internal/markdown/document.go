package markdown

import (
	"io"
)

type Document struct {
	root groupBlock
}

func (d *Document) Append(b ...Block) {
	d.root = append(d.root, b...)
}

func (d *Document) Generate(w io.Writer) error {
	return d.root.Markdown(w)
}
