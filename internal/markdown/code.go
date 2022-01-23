package markdown

func Code(code string) Block {
	return Wrap([]byte("`"), []byte("`"), T(code))
}

func CodeBlock(code string, syntax ...string) Block {
	var prefix []byte
	if len(syntax) != 0 {
		prefix = []byte("```" + syntax[0] + "\n")
	} else {
		prefix = []byte("```\n")
	}
	return Wrap(prefix, []byte("\n```\n"), T(code))
}
