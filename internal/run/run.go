package run

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/codemicro/brainfuck/internal/def"
)

var (
	ErrorNegativePointer = errors.New("interpreter error: pointer cannot be negative")
	// ErrorPointerOutOfRange = errors.New("interpreter error: pointer exceeded maximum value")
	ErrorNoLoops          = errors.New("interpreter error: attempt to end loop without corresponding starter (has this been parsed?)")
	ErrorIllegalCharacter = errors.New("interpreter error: illegal character")
)

func Run(in []byte, input io.Reader, output io.Writer, bufferOutput bool) error {
	var ptr, pc int
	tape := make(memoryTape)

	var loopStarters intStack
	var outputBuffer, inputBuffer []byte

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
			v := tape.Get(ptr)
			if v == 10 {
				if bufferOutput {
					fmt.Fprintln(output, string(outputBuffer))
					outputBuffer = []byte{}
				} else {
					fmt.Fprintln(output)
				}
			} else {
				if bufferOutput {
					outputBuffer = append(outputBuffer, v)
				} else {
					fmt.Fprint(output, string(v))
				}
			}

		case def.SymbolInput:
			// accept one byte of input, storing its value in the byte at the data pointer.

			if len(inputBuffer) == 0 {
				scanner := bufio.NewScanner(input)
				scanner.Scan()
				inputBuffer = append(scanner.Bytes(), 10)
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("runtime error: unable to read from input (%s)", err.Error())
				}
			}

			if len(inputBuffer) != 0 {
				tape.Set(ptr, inputBuffer[0])
				inputBuffer = inputBuffer[1:]
			}

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

	if len(outputBuffer) != 0 {
		fmt.Fprintln(output, string(outputBuffer))
	}

	return nil

}
