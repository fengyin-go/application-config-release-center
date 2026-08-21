package parser

type Parser struct{ buffer []byte }

func (p *Parser) Parse(input string) []byte {
	if cap(p.buffer) < len(input) {
		p.buffer = make([]byte, len(input))
	}
	p.buffer = p.buffer[:len(input)]
	copy(p.buffer, input)
	return p.buffer
}
