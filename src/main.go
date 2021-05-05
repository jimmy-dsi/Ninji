package main

import (
	"fmt"
	. "ninji/lexer"
)

func main() {
	fmt.Println("Hello!")

	lexer := Lexer{}.Init("..\\tests\\1.ninj")
	tokens := lexer.Lex()

	fmt.Println(fmt.Sprint(len(tokens)))

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		fmt.Println("ID:       ", token.ID)
		fmt.Println("RawValue: ", string(token.RawValue))
		fmt.Println("Value:    ", token.Value)
		//fmt.Println("FilePath: ", token.FilePath)
		fmt.Println("Line:     ", token.Line, ":", token.Column)
		//fmt.Println("Column:   ", token.Column)
		fmt.Println("")
	}

	fmt.Println("Done.")
}
