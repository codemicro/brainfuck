package parse

import (
	"errors"

	"github.com/codemicro/brainfuck/internal/def"
)

var (
	ErrorUnbalancedLoop = errors.New("unbalanced loop")
)

func String(in string) ([]byte, error) {
	return parser([]byte(in))
}

func Byte(in []byte) ([]byte, error) {
	return parser(in)
}

func parser(src []byte) ([]byte, error) {

	var o []byte

	for _, bt := range src {

		// ignore chars that are not part of the lang
		switch bt {
		case def.SymbolIncPtr:
		case def.SymbolDecPtr:
		case def.SymbolIncVal:
		case def.SymbolDecVal:
		case def.SymbolOutput:
		case def.SymbolInput:
		case def.SymbolLoopStart:
		case def.SymbolLoopEnd:
		default:
			continue
		}

		o = append(o, bt)
	}

	// remove initial comment loop while also checking for unbalanced loops
	if len(o) > 0 {
		if o[0] == def.SymbolLoopStart {

			var c int
			for c = 1; c < len(o); c += 1 {
				if !hasUnbalancedLoops(o[c:]) {
					break
				}
			}

			o = o[c:]

			if len(o) == 0 {
				return []byte{}, ErrorUnbalancedLoop
			}
		} else {
			if hasUnbalancedLoops(o) {
				return []byte{}, ErrorUnbalancedLoop
			}
		}
	}

	return o, nil
}

func hasUnbalancedLoops(o []byte) bool {
	var c int
	for _, bt := range o {
		if bt == def.SymbolLoopStart {
			c += 1
		} else if bt == def.SymbolLoopEnd {
			c -= 1
		}

		// This will detect a loop that is backwards
		// ie. +]-[ will cause this to fail
		if c < 0 {
			return true
		}

	}
	return c != 0
}
