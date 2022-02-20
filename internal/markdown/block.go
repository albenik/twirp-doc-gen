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

// G renders a group of blocks.
type blockGroup []Block

func G(blocks ...Block) Block {
	if len(blocks) == 0 {
		return emptyBlock{}
	}

	if len(blocks) == 1 {
		return blocks[0]
	}

	return blockGroup(blocks)
}

func (g blockGroup) Markdown(w io.Writer) error {
	suffix := newCatcher(newline)
	mw := io.MultiWriter(w, suffix)

	for _, block := range g {
		// if newline was written by the previous block
		if suffix.Found() {
			// then write second newline to separate blocks as it required by markdown
			if _, err := w.Write(newline); err != nil {
				return err
			}
		}

		if err := block.Markdown(mw); err != nil {
			return err
		}
	}

	return nil
}

type wrapperBlock struct {
	Prefix []byte
	Suffix []byte
	Nested Block
}

// Wrap wraps another block with prefix & suffix.
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

type catcher struct {
	target []byte
	buffer []byte
	buflen int
}

func newCatcher(p []byte) *catcher {
	return &catcher{
		target: p,
		buffer: make([]byte, len(p)),
		buflen: len(p),
	}
}

func (c *catcher) Write(p []byte) (int, error) {
	if l := len(p); l >= c.buflen {
		copy(c.buffer, p[l-c.buflen:])
	} else {
		copy(c.buffer, c.buffer[l:])
		copy(c.buffer[c.buflen-l:], p)
	}

	return len(p), nil
}

func (c *catcher) Found() bool {
	return bytes.Equal(c.buffer, c.target)
}
