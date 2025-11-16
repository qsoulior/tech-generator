package main

import (
	"fmt"

	"github.com/expr-lang/expr"
)

func main() {
	env := map[string]any{
		"foo":    100,
		"bar":    200,
		"foobar": "foo + bar",
		"test":   500,
	}

	program, err := expr.Compile(`foobar + test`)
	if err != nil {
		fmt.Println("1 - ", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Println("2 - ", err)
		return
	}

	fmt.Print(output) // 300
}
