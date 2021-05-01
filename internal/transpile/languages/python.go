package languages

import (
	"bytes"
	"fmt"
)

type Python struct {
	currentTabNumber int
	buf *bytes.Buffer
}

func NewPython() *Python {
	var x []byte
	return &Python{
		buf: bytes.NewBuffer(x),
	}
}

func (p *Python) applyTabs() {
	p.buf.Write(bytes.Repeat([]byte("\t"), p.currentTabNumber))
}

func (p *Python) Header() {
	p.buf.Write([]byte(`input_buffer = ""
def take_input():
    global input_buffer
    if len(input_buffer) == 0:
        input_buffer = (input() + "\n").encode()

    x = input_buffer[0]
    input_buffer = input_buffer[1:]
    return x

def output(x):
    print(chr(x), end="")

memory = {}

def mem_get(ptr):
    return memory.get(ptr, 0)

def mem_set(ptr, val):
    if val == 0:
        del memory[ptr]
    else:
        memory[ptr] = val

def mem_apply_delta(ptr, delta):
    mem_set(ptr, mem_get(ptr) + delta)
`))
}

func (p *Python) MainStart() {
	p.buf.Write([]byte(`def main():
	ptr = 0
`))
	p.currentTabNumber += 1
}

func (p *Python) MainEnd() {
	p.buf.Write([]byte(`
if __name__ == "__main__":
	main()
`))
	p.currentTabNumber -= 1
}

func (p *Python) ValueDelta(delta int) {
	p.applyTabs()
	p.buf.Write([]byte(fmt.Sprintf(`mem_apply_delta(ptr, %d)
`, delta)))
}

func (p *Python) PointerDelta(delta int) {
	p.applyTabs()
	p.buf.Write([]byte(fmt.Sprintf(`ptr += %d
`, delta)))
}

func (p *Python) Output() {
	p.applyTabs()
	p.buf.Write([]byte(`output(mem_get(ptr))
`))
}

func (p *Python) Input() {
	p.applyTabs()
	p.buf.Write([]byte(`mem_set(ptr, take_input())
`))
}

func (p *Python) LoopStart() {
	p.applyTabs()
	p.buf.Write([]byte(`while mem_get(ptr) != 0:
`))
	p.currentTabNumber += 1
}

func (p *Python) LoopEnd() {
	p.currentTabNumber -= 1
}

func (p *Python) Bytes() []byte {
	return p.buf.Bytes()
}