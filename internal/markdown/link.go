package markdown

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var anchroRx = regexp.MustCompile(`[^a-z0-9\-]`)

type linkBlock struct {
	href  string
	label string
}

func (l *linkBlock) Markdown(w io.Writer) error {
	_, err := fmt.Fprintf(w, "[%s](%s)", l.label, l.href)
	return err
}

func LinkToHeader(hdr, label string) Block {
	return &linkBlock{
		href:  "#" + anchroRx.ReplaceAllString(strings.ReplaceAll(strings.ToLower(hdr), " ", "-"), ""),
		label: label,
	}
}

func Link(href, label string) Block {
	return &linkBlock{
		href:  href,
		label: label,
	}
}
