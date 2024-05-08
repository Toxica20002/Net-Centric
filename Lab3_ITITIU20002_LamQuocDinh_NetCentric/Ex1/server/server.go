package main

import (
	// "fmt"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

var NumberOfUser = 0
var UserList = make(map[int]int)

func main() {

	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	request := string(buffer[:])
	parts := strings.Split(request, " ")

	if parts[0] == "login" {
		// write data to response
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(100) + 1
		NumberOfUser++
		response := strconv.Itoa(NumberOfUser) + " "
		UserList[NumberOfUser] = number
		conn.Write([]byte(response))
	} else if parts[0] == "guess" {
		userID, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}

		guess, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Fatal(err)
		}

		number, ok := UserList[userID]
		if !ok {
			fmt.Println("(Server) Error: User not found")
			os.Exit(1)
		}

		if guess == number {
			response := "Correct "
			conn.Write([]byte(response))
		} else if guess > number {
			response := "Lower "
			conn.Write([]byte(response))
		} else {
			response := "Higher "
			conn.Write([]byte(response))
		}

	}

	// close conn
	conn.Close()
}
