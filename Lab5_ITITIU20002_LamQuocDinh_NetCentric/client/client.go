package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type MessageQueue []string

func (mq *MessageQueue) Push(message string) {
	*mq = append(*mq, message)
}

func (mq *MessageQueue) Pop() (string, bool) {
	if len(*mq) == 0 {
		return "", false
	}
	message := (*mq)[0]
	*mq = (*mq)[1:]
	return message, true
}

var messageQueue MessageQueue

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
		fmt.Print("Enter recipient (@username or @all) or Check message(@message): ")
		recipient, _ := reader.ReadString('\n')
		recipient = strings.TrimSpace(recipient)

		if recipient != "@message" {
			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			fullMessage := recipient + " " + username + " " + message
			_, err = conn.Write([]byte(fullMessage))
			checkError(err)
		} else {

			//Print messages from the server
			fmt.Println("\n***Messages from the server:***")

			for {
				message, ok := messageQueue.Pop()
				if !ok {
					break
				}
				fmt.Println(message)
			}

			fmt.Println("***End of messages from the server***\n")

		}

	}

	os.Exit(0)
}

func readMessages(conn *net.UDPConn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		checkError(err)
		message := string(buf[:n])
		messageQueue.Push(message)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
