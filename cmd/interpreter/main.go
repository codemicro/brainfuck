package main

import (
	"fmt"

	"github.com/codemicro/brainfuck/internal/parse"
	"github.com/codemicro/brainfuck/internal/run"
)

func main() {
	parsed, err := parse.String(`++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.`)

	if err != nil {
		panic(err)
	}

	for _, x := range parsed {
		fmt.Print(string(x))
	}
	fmt.Println()

	err = run.Run(parsed)

	if err != nil {
		panic(err)
	}

}
