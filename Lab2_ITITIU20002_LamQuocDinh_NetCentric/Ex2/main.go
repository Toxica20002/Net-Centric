package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MaxCapacity   = 30
	TotalStudents = 100
)

var startTime = time.Now()

func student(wg *sync.WaitGroup, lib chan int, id int) {
	defer wg.Done()

	// Generate a random reading time from 1 to 4 hours.
	rand.Seed(time.Now().UnixNano())
	readingTime := time.Duration(rand.Intn(4) + 1)

	// Try to enter the library.
	lib <- id
	fmt.Printf("Time %d: Student %d start reading at the lib\n", int(time.Since(startTime).Seconds()), id)

	// Simulate the reading time.
	time.Sleep(readingTime * time.Second)

	// Leave the library.
	<-lib
	fmt.Printf("Time %d: Student %d is leaving. Spent %d hours reading\n", int(time.Since(startTime).Seconds()), id, readingTime)
}

func main() {
	lib := make(chan int, MaxCapacity)
	var wg sync.WaitGroup

	for id := 1; id <= TotalStudents; id++ {
		wg.Add(1)
		go student(&wg, lib, id)
	}

	wg.Wait()

	fmt.Printf("Time %d: No more students. Let's call it a day\n", int(time.Since(startTime).Seconds()))
	fmt.Printf("The library needs to be open for %d hours\n", int(time.Since(startTime).Seconds()))
}
