package markdown_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	md "github.com/albenik/twirp-doc-gen/internal/markdown"
)

type testCase struct {
	Name   string
	Block  md.Block
	Result string
}

func runTestCases(t *testing.T, cases []*testCase) {
	t.Helper()

	for _, c := range cases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer(nil)
			require.NoError(t, c.Block.Markdown(buf))
			require.Equal(t, c.Result, buf.String())
		})
	}
}
