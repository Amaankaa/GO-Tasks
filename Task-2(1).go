package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a string: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	slice := strings.Fields(input)

	dict := make(map[string]int)

	for _, val := range slice {
		dict[strings.ToLower(val)]++
	}

	fmt.Print(dict)
}