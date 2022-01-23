package markdown

import (
	"bytes"
	"io"
)

var newline = []byte("\n")

type Block interface {
	Markdown(io.Writer) error
}

type emptyBlock struct{}

func (emptyBlock) Markdown(io.Writer) error {
	return nil
}

// G renders a group of blocks
type groupBlock []Block

func G(blocks ...Block) Block {
	if len(blocks) == 0 {
		return emptyBlock{}
	}

	if len(blocks) == 1 {
		return blocks[0]
	}

	return groupBlock(blocks)
}

func (b groupBlock) Markdown(w io.Writer) (err error) {
	suffix := newSuffixCatcher(newline)
	mw := io.MultiWriter(w, suffix)

	for _, bb := range b {
		if suffix.Found() {
			_, err = w.Write(newline)
		}
		if err = bb.Markdown(mw); err != nil {
			return
		}
	}

	return
}

type wrapperBlock struct {
	Prefix []byte
	Suffix []byte
	Nested Block
}

// Wrap wraps another block with prefix & suffix
func Wrap(prefix, suffix []byte, blocks ...Block) Block {
	return &wrapperBlock{
		Prefix: prefix,
		Suffix: suffix,
		Nested: G(blocks...),
	}
}

func (b *wrapperBlock) Markdown(w io.Writer) error {
	if b.Prefix != nil {
		if _, err := w.Write(b.Prefix); err != nil {
			return err
		}
	}

	if err := b.Nested.Markdown(w); err != nil {
		return err
	}

	if b.Suffix != nil {
		if _, err := w.Write(b.Suffix); err != nil {
			return err
		}
	}

	return nil
}

type suffixCatcher struct {
	suffix []byte
	buffer []byte
}

func newSuffixCatcher(s []byte) *suffixCatcher {
	return &suffixCatcher{
		suffix: s,
		buffer: make([]byte, len(s)),
	}
}

func (c *suffixCatcher) Write(b []byte) (int, error) {
	if l := len(b); l >= 2 {
		copy(c.buffer, b[l-2:])
	} else {
		copy(c.buffer, c.buffer[l:])
		copy(c.buffer[len(c.buffer)-1:], b)
	}
	return len(b), nil
}

func (c *suffixCatcher) Found() bool {
	return bytes.Equal(c.buffer, c.suffix)
}
