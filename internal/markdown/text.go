package markdown

import (
	"io"
)

var (
	asterisk   = []byte("*")
	asteriskX2 = []byte("**")
	hr         = []byte("---\n")
)

type textBlock []byte

// T renders a plain text block.
func T(s string) Block {
	return textBlock(s)
}

func (t textBlock) Markdown(w io.Writer) error {
	_, err := w.Write(t)
	return err
}

// I renders an italic text block.
func I(blocks ...Block) Block {
	return Wrap(asterisk, asterisk, blocks...)
}

func TI(s string) Block {
	return I(T(s))
}

// B renders a bold text block.
func B(blocks ...Block) Block {
	return Wrap(asteriskX2, asteriskX2, blocks...)
}

func TB(s string) Block {
	return B(T(s))
}

// P renders a paragraph.
func P(blocks ...Block) Block {
	return Wrap(nil, newline, blocks...)
}

func TP(s string) Block {
	return P(T(s))
}

type lineBlock struct{}

func Line() Block {
	return lineBlock{}
}

func (lineBlock) Markdown(w io.Writer) error {
	_, err := w.Write(hr)
	return err
}
