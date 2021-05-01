package main

import (
	"errors"
	"fmt"
	"github.com/codemicro/brainfuck/internal/parse"
	"github.com/codemicro/brainfuck/internal/transpile"
	"github.com/codemicro/brainfuck/internal/transpile/languages"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
)

const inputString = "command"
const outputFile = "outputFile"
const outputLang = "outputLanguage"

func main() {
	app := &cli.App{
		Name:  "transpile",
		Usage: "transpile a brainfuck program",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: inputString,
				Aliases: []string{"c"},
				Usage: "executes an input string. Takes precedence over any file",
			},
			&cli.StringFlag{
				Name: outputFile,
				Aliases: []string{"o"},
				Required: true,
				Usage: "output filename",
			},
			&cli.StringFlag{
				Name: outputLang,
				Aliases: []string{"l"},
				DefaultText: "go",
				Usage: "language for transpiled output",
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

			ol := c.String(outputLang)
			if ol == "" {
				ol = "go"
			}
			lang, found := languages.Index[ol]
			if !found {
				return errors.New("unknown output language")
			}

			transpiled, err := transpile.Transpile(parsed, lang)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(c.String(outputFile), transpiled, 0644)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
