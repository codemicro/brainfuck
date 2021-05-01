package languages

import (
	"bytes"
	"fmt"
)

type Go struct {
	buf *bytes.Buffer
}

func NewGo() *Go {
	var ob []byte
	return &Go{
		buf: bytes.NewBuffer(ob),
	}
}

func (g *Go) Header() {
	g.buf.Write([]byte(`package main

import (
	"bufio"
	"fmt"
	"os"
)

var inputBuffer []byte
func takeInput() (byte, error) {
	if len(inputBuffer) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		inputBuffer = append(scanner.Bytes(), 10)
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("runtime error: unable to read from input (%s)", err.Error())
		}
	}

	x := inputBuffer[0]
	inputBuffer = inputBuffer[1:]

	return x, nil
}

func output(x byte) {
	_, _ = os.Stdout.Write([]byte{x})
}

type mem map[int]byte

func (t mem) Get(ptr int) byte {
	v, ok := t[ptr]
	if !ok {
		return 0
	}
	return v
}

func (t mem) Set(ptr int, val byte) {
	if val == 0 {
		delete(t, ptr)
	} else {
		t[ptr] = val
	}
}

func (t mem) ApplyDelta(ptr, delta int) {
	t.Set(ptr, byte(int(t.Get(ptr))+delta))
}
`))
}

func (g *Go) MainStart() {
	g.buf.Write([]byte(`func main() {
pm := make(mem)
var ptr int
`))
}

func (g *Go) MainEnd() {
	g.buf.Write([]byte(`}
`))
}

func (g *Go) ValueDelta(delta int) {
	g.buf.Write([]byte(fmt.Sprintf(`pm.ApplyDelta(ptr, %d)
`, delta)))
}

func (g *Go) PointerDelta(delta int) {
	g.buf.Write([]byte(fmt.Sprintf(`ptr += %d
`, delta)))
}

func (g *Go) Output() {
	g.buf.Write([]byte(`output(pm.Get(ptr))
`))
}

func (g *Go) Input() {
	g.buf.Write([]byte(`pm.Set(ptr, takeInput())
`))
}

func (g *Go) LoopStart() {
	g.buf.Write([]byte(`for pm.Get(ptr) != 0 {
`))
}

func (g *Go) LoopEnd() {
	g.buf.Write([]byte(`}
`))
}

func (g *Go) Bytes() []byte {
	return g.buf.Bytes()
}