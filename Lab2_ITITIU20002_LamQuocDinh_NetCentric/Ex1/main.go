package main

import (
	"fmt"
)

func countLetter(s string, chanCount chan int, char rune) {
	count := 0
	for _, c := range s {
		if c == char {
			count++
		}
	}
	chanCount <- count
	
}

func main() {
	chanCount := make(chan int)
	frequency := make(map[rune]int)
   
	fmt.Print("Enter a string: ")
	var s string
	fmt.Scanln(&s)


	for _, c := range s {
		go countLetter(s, chanCount, c)
		frequency[c] = <-chanCount
	}
	
	for k, v := range frequency {
		if v > 0 {
			if k == ' '{
				fmt.Println("(blank)", ":", v)
			} else {
				fmt.Println(string(k), ":", v)
			}
				
		}
	}
}