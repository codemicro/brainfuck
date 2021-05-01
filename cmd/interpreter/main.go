package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codemicro/brainfuck/internal/parse"
	"github.com/codemicro/brainfuck/internal/run"
	"github.com/urfave/cli/v2"
)

// http://www.hevanet.com/cristofd/brainfuck/

const bufferOutput = "disableOutputBuffer"
const inputString = "command"

func main() {
	app := &cli.App{
		Name:  "run",
		Usage: "run a brainfuck program",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: inputString,
				Aliases: []string{"c"},
				Usage: "executes an input string. Takes precendence over any file",
			},
			&cli.BoolFlag{
				Name: bufferOutput,
				Aliases: []string{"b"},
				Usage: "disables output buffering",
			},
		},
		Action: func(c *cli.Context) error {

			var parsed []byte

			if inputString := c.String(inputString); inputString != "" {
				var err error
				parsed, err = parse.String(inputString)
				if err != nil {
					return err
				}
			} else if inputFile := c.Args().Get(0); inputFile != "" {
				x, err := ioutil.ReadFile(inputFile)
				if err != nil {
					return fmt.Errorf("unable to read input file %s: %s", inputFile, err.Error())
				}

				parsed, err = parse.Bytes(x)
				if err != nil {
					return err
				}
			}

			return run.Run(parsed, os.Stdin, os.Stdout, !c.Bool(bufferOutput))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
