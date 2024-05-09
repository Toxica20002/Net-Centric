package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type ClientInfo struct {
	IP   string
	Name string
}

var clients = make(map[string]ClientInfo)

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	message := string(buf[:n])
	fmt.Println("Received: ", message)
	//fmt.Println("Received: ", message)
	if strings.HasPrefix(message, "@") {
		parts := strings.SplitN(message, " ", 3)
		command := parts[0][1:] // remove the '@'
		if len(parts) > 1 {
			msg := parts[2]
			sender := parts[1]
			if command == "all" {
				for addr, client := range clients {
					if client.Name == sender {
						continue
					}
					msg = "(Public) " + sender + ": " + msg
					clientAddr, _ := net.ResolveUDPAddr("udp", addr) // assuming clients are listening on port 1200
					conn.WriteToUDP([]byte(msg), clientAddr)
				}
			} else {
				msg = "(Private) " + sender + ": " + msg
				name := command
				name = strings.TrimSpace(name)
				flag := false
				for addr, client := range clients {
					if client.Name == name {
						clientAddr, _ := net.ResolveUDPAddr("udp", addr) // assuming clients are listening on port 1200
						conn.WriteToUDP([]byte(msg), clientAddr)
						flag = true
					}
				}
				if !flag {
					errMsg := "(Public) Server: The user " + command + " does not exist."
					conn.WriteToUDP([]byte(errMsg), addr)
				}
			}
		}
	} else if client, ok := clients[addr.String()]; ok {
		clients[addr.String()] = client
	} else {
		// fmt.Println("New client registered: " + message + " from " + addr.String())
		message = strings.TrimSpace(message)
		clients[addr.String()] = ClientInfo{IP: addr.IP.String(), Name: message}

		for addr, _ := range clients {
			msg := "(Public) Server: User " + message + " has joined the chat."
			clientAddr, _ := net.ResolveUDPAddr("udp", addr) // assuming clients are listening on port 1200
			conn.WriteToUDP([]byte(msg), clientAddr)
		}
		fmt.Println("User " + message + " has joined the chat.")
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
