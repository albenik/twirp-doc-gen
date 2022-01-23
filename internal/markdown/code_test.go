package markdown_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/albenik-go/twirp-doc-gen/internal/markdown"
)

func TestInlineCode_Markdown(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	require.NoError(t, markdown.Code("inline code").Markdown(buf))
	require.Equal(t, "`inline code`", buf.String())
}

func TestCodeBlock_Markdown(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	require.NoError(t, markdown.CodeBlock("{\n  foo = bar\n}").Markdown(buf))
	require.Equal(t, "```\n{\n  foo = bar\n}\n```\n", buf.String())
}
