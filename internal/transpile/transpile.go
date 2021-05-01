package transpile

import (
	"github.com/codemicro/brainfuck/internal/def"
	"github.com/codemicro/brainfuck/internal/transpile/languages"
)

func countPrefix(dat []byte, startPos int, symbol byte) int {
	var n int
	if startPos >= len(dat) {
		return 0
	}
	for _, char := range dat[startPos:] {
		if char != symbol {
			break
		}
		n += 1
	}
	return n
}

func Transpile(in []byte, transpiler languages.Language) ([]byte, error) {

	transpiler.Header()
	transpiler.MainStart()

	var ptr int

	for ptr < len(in) {

		cir := in[ptr]

		switch cir {
		case def.SymbolIncPtr:
			// increment the data pointer.
			n := countPrefix(in, ptr+1, def.SymbolIncPtr) // plus one to run from the next item in the input
			transpiler.PointerDelta(n + 1) // plus one to include the symbol that the switch statement has been triggered with
			ptr += n
		case def.SymbolDecPtr:
			// decrement the data pointer.
			n := countPrefix(in, ptr+1, def.SymbolDecPtr)
			transpiler.PointerDelta(-(n + 1))
			ptr += n
		case def.SymbolIncVal:
			// increment the byte at the data pointer.
			n := countPrefix(in, ptr+1, def.SymbolIncVal)
			transpiler.ValueDelta(n + 1)
			ptr += n
		case def.SymbolDecVal:
			// decrement the byte at the data pointer.
			n := countPrefix(in, ptr+1, def.SymbolDecVal)
			transpiler.ValueDelta(-(n + 1))
			ptr += n
		case def.SymbolOutput:
			// output the byte at the data pointer.
			transpiler.Output()
		case def.SymbolInput:
			// accept one byte of input, storing its value in the byte at the data pointer.
			transpiler.Input()
		case def.SymbolLoopStart:
			// if the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ] command.
			transpiler.LoopStart()
		case def.SymbolLoopEnd:
			// if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command.
			transpiler.LoopEnd()
		default:
			// TODO: handle this
		}

		ptr += 1
	}

	transpiler.MainEnd()

	return transpiler.Bytes(), nil
}
