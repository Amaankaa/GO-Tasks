package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName) // removes the trailing newline

	var numberOfSubjects int
	fmt.Print("Enter the number of subjects you took: ")
	fmt.Scanln(&numberOfSubjects)

	subjectToGrade := make(map[string]int)

	total := 0
	for i:=0; i<numberOfSubjects; i++{
		var subject string
		fmt.Printf("Enter the name of subject(%v): ", i + 1)
		fmt.Scanln(&subject)
		var grade int
		fmt.Printf("Enter the grade of subject(%v): ", i + 1)
		fmt.Scanln(&grade)

		for grade < 0 || grade > 100 {
			fmt.Print("Invalid grade. Please enter a valid grade: ")
			fmt.Scanln(&grade)
		}

		subjectToGrade[subject] = grade
		total += subjectToGrade[subject]
	}

	fmt.Println(fullName)
	for key, val := range subjectToGrade {
		fmt.Printf("Subject: %v | Grade: %v", key, val)
		fmt.Print("\n")
	}
	fmt.Println("Average Grade", calculateAverage(total, numberOfSubjects))
}

func calculateAverage(total, numberOfSubjects int) float64 {
	avg := float64(total) / float64(numberOfSubjects)
	return float64(int(avg*100)) / 100
}