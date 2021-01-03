package run

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/codemicro/brainfuck/internal/def"
)

var (
	ErrorNegativePointer = errors.New("runtime error: pointer cannot be negative")
	// ErrorPointerOutOfRange = errors.New("runtime error: pointer exceeded maximum value")
	ErrorEOF              = errors.New("runtime error: got end-of-file")
	ErrorNoLoops          = errors.New("runtime error: attempt to end loop without corresponding starter (has this been parsed?)")
	ErrorIllegalCharacter = errors.New("runtime error: illegal character")
)

func Run(in []byte) error {
	var ptr, pc int
	tape := make(memoryTape)

	var loopStarters stack

	for pc < len(in) {

		cir := in[pc]

		switch cir {
		case def.SymbolIncPtr:
			// increment the data pointer.
			ptr += 1
		case def.SymbolDecPtr:
			// decrement the data pointer.
			ptr -= 1
		case def.SymbolIncVal:
			// increment the byte at the data pointer.
			tape.ApplyDelta(ptr, 1)
		case def.SymbolDecVal:
			// decrement the byte at the data pointer.
			tape.ApplyDelta(ptr, -1)
		case def.SymbolOutput:
			// output the byte at the data pointer.
			fmt.Print(string(tape.Get(ptr)))
		case def.SymbolInput:
			// accept one byte of input, storing its value in the byte at the data pointer.
			inp := make([]byte, 1)
			_, err := os.Stdin.Read(inp)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return ErrorEOF
				}
				return fmt.Errorf("runtime error: unable to read from input (%s)", err.Error())
			}
			tape.Set(ptr, inp[0])
		case def.SymbolLoopStart:
			// if the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ] command.
			if tape.Get(ptr) == 0 {
				var depth int
				nx := pc + 1
				for {
					if in[nx] == def.SymbolLoopEnd {
						if depth == 0 {
							pc = nx
							break
						} else {
							depth -= 1
						}
					} else if in[nx] == def.SymbolLoopStart {
						depth += 1
					}
					nx += 1
				}
			} else {
				loopStarters.Push(pc)
			}

		case def.SymbolLoopEnd:
			// if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command.
			if tape.Get(ptr) != 0 {
				v, err := loopStarters.Top()
				if err != nil {
					return ErrorNoLoops
				}
				pc = v
			} else {
				_, err := loopStarters.Pop()
				if err != nil {
					return ErrorNoLoops
				}
			}
		default:
			return ErrorIllegalCharacter
		}

		if ptr < 0 {
			return ErrorNegativePointer
		}

		pc += 1

	}

	return nil

}
