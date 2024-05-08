package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Stack []int

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(value int) {
	*s = append(*s, value)
}

func (s *Stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

func solve(a string, b string) int {
	n := len(a)
	ans := 0

	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			ans++
		}
	}
	return ans
}

func randomLength() int {
	return rand.Intn(100)
}

func randomADN(ADN []string) string {
	return ADN[rand.Intn(4)]
}

func question1() {
	fmt.Println("Question 1:")
	turn := 1000
	for i := 0; i < turn; i++ {
		fmt.Printf("Pair #%d: \n", i+1)
		ADN := []string{"A", "C", "G", "T"}
		n := randomLength()
		a := ""
		b := ""
		for i := 0; i < n; i++ {
			a += randomADN(ADN)
			b += randomADN(ADN)
		}
		fmt.Println(a)
		fmt.Println(b)
		fmt.Printf("The Hamming Distance is %d \n", solve(a, b))
		fmt.Println()
	}
}

func randomWord(length int, alpha string) string {
	word := ""
	for i := 0; i < length; i++ {
		word += string(alpha[rand.Intn(len(alpha))])
	}
	word = strings.TrimLeft(word, " ")
	word = strings.TrimRight(word, " ")
	return word
}

func question2() {
	var word string
	var length int

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\nQuestion 2:")
	fmt.Printf("Do you want to enter a word? (Y/N): ")
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	word = scanner.Text()
	word = strings.ToUpper(word)
	if word == "Y" {
		fmt.Printf("Enter a word: ")
		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			fmt.Println(err)
			return
		}
		word = scanner.Text()
		fmt.Printf("The word is %s\n", word)
		word = strings.ToUpper(word)
		length = len(word)
	} else {
		alpha := " ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		length = randomLength()
		word = randomWord(length, alpha)
		fmt.Printf("The word is %s\n", word)
	}

	dict := map[string]int{
		" ": 0,
		"A": 1, "E": 1, "I": 1, "O": 1, "U": 1, "L": 1, "N": 1, "R": 1, "S": 1, "T": 1,
		"D": 2, "G": 2,
		"B": 3, "C": 3, "M": 3, "P": 3,
		"F": 4, "H": 4, "V": 4, "W": 4, "Y": 4,
		"K": 5,
		"J": 8, "X": 8,
		"Q": 10, "Z": 10,
	}
	score := 0
	length = len(word)
	for i := 0; i < length; i++ {
		score += dict[string(word[i])]
	}
	fmt.Printf("The Scrabble score is %d\n", score)
}

func randomCreditCardNumber() []int {
	creditCardNumber := make([]int, 16)
	for i := 0; i < 16; i++ {
		creditCardNumber[i] = rand.Intn(10)
	}
	return creditCardNumber
}

func question3() {
	var creditCardNumber []int
	fmt.Println("\nQuestion 3:")

	fmt.Printf("Do you want to enter a credit card number? (Y/N): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	answer := scanner.Text()
	answer = strings.ToUpper(answer)
	if answer == "Y" {
		fmt.Printf("Enter with format XXXX XXXX XXXX XXXX: ")
		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			fmt.Println(err)
			return
		}
		creditCardNumber = make([]int, 16)
		creditCard := scanner.Text()
		creditCard = strings.ReplaceAll(creditCard, " ", "")
		for i := 0; i < 16; i++ {
			creditCardNumber[i] = int(creditCard[i] - '0')
		}
	} else {
		creditCardNumber = randomCreditCardNumber()
	}

	fmt.Printf("The credit card number is ")
	for i := 0; i < 16; i++ {
		fmt.Printf("%d", creditCardNumber[i])
		if i%4 == 3 {
			fmt.Printf(" ")
		}
	}
	for i := 0; i < 16; i++ {
		if i%2 == 0 {
			creditCardNumber[i] *= 2
			if creditCardNumber[i] > 9 {
				creditCardNumber[i] -= 9
			}
		}
	}
	fmt.Printf("\nThe credit card number after Luhn algorithm is ")
	for i := 0; i < 16; i++ {
		fmt.Printf("%d", creditCardNumber[i])
		if i%4 == 3 {
			fmt.Printf(" ")
		}
	}
	sum := 0
	for i := 0; i < 16; i++ {
		sum += creditCardNumber[i]
	}
	if sum%10 == 0 {
		fmt.Println("\nThe credit card number is valid")
	} else {
		fmt.Println("\nThe credit card number is invalid")
	}

}

func question4() {
	fmt.Println("\nQuestion 4:")
	rows := 20
	cols := 25
	minefield := make([][]string, rows)
	for i := 0; i < rows; i++ {
		minefield[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			state := rand.Intn(2)
			if state == 0 {
				minefield[i][j] = "."
			} else {
				minefield[i][j] = "*"
			}
		}
	}

	fmt.Printf("The minefield is:\n")
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf("%s", minefield[i][j])
		}
		fmt.Println()
	}

	fmt.Printf("The minefield after processing is:\n")
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if minefield[i][j] == "." {
				countMine := 0
				for di := -1; di <= 1; di++ {
					for dj := -1; dj <= 1; dj++ {
						if i+di >= 0 && i+di < rows && j+dj >= 0 && j+dj < cols {
							if minefield[i+di][j+dj] == "*" {
								countMine++
							}
						}
					}
				}
				if countMine == 0 {
					fmt.Printf(".")
				} else {
					fmt.Printf("%d", countMine)
				}
			} else {
				fmt.Printf("*")
			}
		}
		fmt.Println()
	}

}

func question5() {
	var word string
	fmt.Printf("\nQuestion 5:\n")
	fmt.Printf("Do you want to enter a word? (Y/N): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	answer := scanner.Text()
	answer = strings.ToUpper(answer)
	if answer == "Y" {
		fmt.Printf("Enter a word: ")
		scanner.Scan()
		err = scanner.Err()
		if err != nil {
			fmt.Println(err)
			return
		}
		word = scanner.Text()
		fmt.Printf("The word is %s\n", word)
	} else {
		alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ.[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()[]{}()"
		length := randomLength()
		word = randomWord(length, alpha)
		fmt.Printf("The word is %s\n", word)
	}

	var myStack Stack
	for i := 0; i < len(word); i++ {
		if word[i] == '(' || word[i] == '[' || word[i] == '{' {
			myStack.Push(int(word[i]))
		} else if word[i] == ')' {
			element, ok := myStack.Pop()
			if !ok || element != int('(') {
				fmt.Printf("Incorrect\n")
				return
			}
		} else if word[i] == ']' {
			element, ok := myStack.Pop()
			if !ok || element != int('[') {
				fmt.Printf("Incorrect\n")
				return
			}
		} else if word[i] == '}' {
			element, ok := myStack.Pop()
			if !ok || element != int('{') {
				fmt.Printf("Incorrect\n")
				return
			}
		}
	}
	if myStack.IsEmpty() {
		fmt.Printf("Correct\n")
	} else {
		fmt.Printf("Incorrect\n")
	}

}

func main() {
	question1()
	question2()
	question3()
	question4()
	question5()
}
