package repl

import (
	"bufio"
	"fmt"
	"io"
	"rog/lexer"
	// "rog/token"
	"rog/parser"
)

const PROMPT = ">> "
const ROG_ICON = ` ______     ______     ______    
/\  == \   /\  __ \   /\  ___\   
\ \  __<   \ \ \/\ \  \ \ \__ \  
 \ \_\ \_\  \ \_____\  \ \_____\ 
  \/_/ /_/   \/_____/   \/_____/                                

`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	io.WriteString(out, ROG_ICON)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return 
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}