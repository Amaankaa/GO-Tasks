package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a string: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Remove all non-letter characters
	var cleaned strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) {
			cleaned.WriteRune(unicode.ToLower(r))
		}
	}
	cleanedStr := cleaned.String()

	left, right := 0, len(cleanedStr)-1
	flag := true

	for left < right {
		if cleanedStr[left] != cleanedStr[right] {
			flag = false
			break
		}
		left++
		right--
	}

	fmt.Println(flag)
}
