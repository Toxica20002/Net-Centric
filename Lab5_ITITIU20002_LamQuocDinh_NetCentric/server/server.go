package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type ClientInfo struct {
	IP            string
	Name          string
	LastHeartbeat time.Time
}

var clients = make(map[string]ClientInfo)
var mutex = &sync.Mutex{}

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	go checkHeartbeats()
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

	mutex.Lock()
	//fmt.Println("Received: ", message)
	if strings.HasPrefix(message, "@") {
		parts := strings.SplitN(message, " ", 2)
		command := parts[0][1:] // remove the '@'
		fmt.Println("message: ", message)
		fmt.Println("command: ", command)
		if len(parts) > 1 {
			msg := parts[1]
			if command == "all" {
				// fmt.Println("Broadcasting message: ", msg)
				// fmt.Println("Clients: ", clients)
				for key, client := range clients {
					fmt.Println("Client: ", client)
					clientAddr, _ := net.ResolveUDPAddr("udp", key) // assuming clients are listening on port 1200
					conn.WriteToUDP([]byte(msg), clientAddr)
					fmt.Println("ClientAddr: ", clientAddr)
					fmt.Println("Key: ", key)
				}
			} else {
				// if client, ok := clients[command]; ok {
				// 	clientAddr, _ := net.ResolveUDPAddr("udp", client.IP+":1200") // assuming clients are listening on port 1200
				// 	conn.WriteToUDP([]byte(msg), clientAddr)
				// } else {
				// 	errMsg := "The user " + command + " does not exist."
				// 	conn.WriteToUDP([]byte(errMsg), addr)
				// }
				name := command
				name = strings.TrimSpace(name)
				flag := false
				for key, client := range clients {
					if client.Name == name {
						clientAddr, _ := net.ResolveUDPAddr("udp", key) // assuming clients are listening on port 1200
						conn.WriteToUDP([]byte(msg), clientAddr)
						flag = true
					}
				}
				if !flag {
					errMsg := "The user " + command + " does not exist."
					conn.WriteToUDP([]byte(errMsg), addr)
				}
			}
		}
	} else if client, ok := clients[addr.String()]; ok {
		client.LastHeartbeat = time.Now()
		clients[addr.String()] = client
	} else {
		// fmt.Println("New client registered: " + message + " from " + addr.String())
		message = strings.TrimSpace(message)
		clients[addr.String()] = ClientInfo{IP: addr.IP.String(), Name: message, LastHeartbeat: time.Now()}
		fmt.Println("New client registered: " + message + " from " + addr.String())
	}
	mutex.Unlock()
}

func checkHeartbeats() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for addr, client := range clients {
			if time.Since(client.LastHeartbeat) > 10*time.Second {
				delete(clients, addr)
			}
		}
		mutex.Unlock()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
