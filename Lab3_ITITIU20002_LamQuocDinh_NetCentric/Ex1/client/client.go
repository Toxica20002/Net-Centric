package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

var UserID = 0

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("(Client) Error: ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	handleLogin(tcpServer)
	flag := false
	for !flag {
		flag = handleGuess(tcpServer)
	}
}

func handleLogin(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("(Client) Error: Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("login "))
	if err != nil {
		println("(Client) Error: Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("(Client) Error: Read data failed:", err.Error())
		os.Exit(1)
	}

	UserID, err = strconv.Atoi(string(received[0]))
	if err != nil {
		println("(Client) Error: Convert to int failed:", err.Error())
		os.Exit(1)
	}
	conn.Close()
}

func handleGuess(tcpServer *net.TCPAddr) bool {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("(Client) Error: Dial failed:", err.Error())
		os.Exit(1)
	}

	fmt.Print("(Client) Enter your guess: ")
	var guess int
	fmt.Scanln(&guess)
	_, err = conn.Write([]byte("guess " + strconv.Itoa(UserID) + " " + strconv.Itoa(guess) + " "))
	if err != nil {
		println("(Client) Error: Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("(Client) Error: Read data failed:", err.Error())
		os.Exit(1)
	}

	response := string(received[:])

	parts := strings.Split(response, " ")
	fmt.Printf("(Server) Response: %v\n", parts[0])

	if parts[0] == "Correct" {
		conn.Close()
		return true
	}

	conn.Close()
	return false
}
