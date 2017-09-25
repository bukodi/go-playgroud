package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter your text: ")
	scanner.Scan()
	text = scanner.Text()
	fmt.Println("Your text was: ", text)
}
