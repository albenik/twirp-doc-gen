package markdown

var (
	h1 = []byte("# ")
	h2 = []byte("## ")
	h3 = []byte("### ")
	h4 = []byte("#### ")
	h5 = []byte("##### ")
	h6 = []byte("###### ")
)

func H1(blocks ...Block) Block {
	return h(h1, blocks)
}

func TH1(s string) Block {
	return H1(T(s))
}

func H2(blocks ...Block) Block {
	return h(h2, blocks)
}

func TH2(s string) Block {
	return H2(T(s))
}

func H3(blocks ...Block) Block {
	return h(h3, blocks)
}

func TH3(s string) Block {
	return H3(T(s))
}

func H4(blocks ...Block) Block {
	return h(h4, blocks)
}

func TH4(s string) Block {
	return H4(T(s))
}

func H5(blocks ...Block) Block {
	return h(h5, blocks)
}

func TH5(s string) Block {
	return H5(T(s))
}

func H6(blocks ...Block) Block {
	return h(h6, blocks)
}

func TH6(s string) Block {
	return H6(T(s))
}

func h(prefix []byte, blocks []Block) Block {
	return Wrap(prefix, newline, blocks...)
}
