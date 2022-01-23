package markdown

import (
	"fmt"
	"io"
)

type listBlock struct {
	itemPrefix []byte
	items      []Block
}

func (l *listBlock) Markdown(w io.Writer) error {
	for i, item := range l.items {
		if l.itemPrefix == nil {
			if _, err := fmt.Fprintf(w, "%d. ", i+1); err != nil {
				return err
			}
		} else {
			if _, err := w.Write(l.itemPrefix); err != nil {
				return err
			}
		}
		if err := item.Markdown(w); err != nil {
			return err
		}
		if _, err := w.Write(newline); err != nil {
			return err
		}
	}
	return nil
}

func UL(block ...Block) Block {
	return &listBlock{
		itemPrefix: []byte("* "),
		items:      block,
	}
}

func OL(block ...Block) Block {
	return &listBlock{
		itemPrefix: nil,
		items:      block,
	}
}
