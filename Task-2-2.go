package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main(){
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a string: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	left, right := 0, len(input) - 1
	flag := true

	for left < right {
		for (left < right) && !unicode.IsLetter(rune(input[left])) { left++ }
		for (left < right) && !unicode.IsLetter(rune(input[right])) { right-- }

		if !(strings.EqualFold(strings.ToLower(string(input[left])),strings.ToLower(string(input[right])))) {
			flag = false
		}
		left++
		right--
	}

	fmt.Println(flag)
}

