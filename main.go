package main

import (
	"dczombera/monkey_language_interpreter/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the simplistic Monkey programming language!\n",
		user.Username)
	fmt.Printf("Type in any command of the language\n")
	repl.Start(os.Stdin, os.Stdout)
}
