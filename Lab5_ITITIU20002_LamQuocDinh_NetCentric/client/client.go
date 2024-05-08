package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port username", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	username := os.Args[2]
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	// Register the client with the server
	_, err = conn.Write([]byte(username))
	checkError(err)

	go readMessages(conn)

	// Read messages from stdin and send them to the server
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter recipient (@username or @all): ")
		recipient, _ := reader.ReadString('\n')
		recipient = strings.TrimSpace(recipient)

		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		fullMessage := recipient + " " + message
		_, err = conn.Write([]byte(fullMessage))
		checkError(err)

		if message == "logout" {
			break
		}
	}

	os.Exit(0)
}

func readMessages(conn *net.UDPConn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		checkError(err)
		fmt.Println(string(buf[0:n]))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
