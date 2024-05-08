/**
 * Author: Yin Lin
 * Client side of the guessing game
**/

package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
	handleListFiles(tcpServer)
	flag := false
	for !flag {
		flag = handleDownloadFile(tcpServer)
	}
}

func handleListFiles(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("(Client) Error: Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("List "))
	if err != nil {
		println("(Client) Error: Write data failed:", err.Error())
		os.Exit(1)
	}

	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("(Client) Error: Read data failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("(Client) Files available for download:")
	parts := strings.Split(string(received), " ")
	for i := 0; i < len(parts)-1; i++ {
		fmt.Printf("(%d) %s\n", i+1, parts[i])
	}
	conn.Close()
}


func handleLogin(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("(Client) Error: Dial failed:", err.Error())
		os.Exit(1)
	}

	fmt.Print("(Client) Enter your username: ")
	var username string
	fmt.Scanln(&username)
	fmt.Print("(Client) Enter your password: ")
	var password string
	fmt.Scanln(&password)

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	base64Hash := base64.URLEncoding.EncodeToString(hash)
	_, err = conn.Write([]byte("login " + username + " " + base64Hash + " "))
	if err != nil {
		println("(Client) Error: Write data failed:", err.Error())
		os.Exit(1)
	}

	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("(Client) Error: Read data failed:", err.Error())
		os.Exit(1)
	}
	parts := strings.Split(string(received), " ")
	if parts[0] == "Invalid" {
		fmt.Println("(Server) Response: Invalid username or password")
		os.Exit(1)
	}

	UserID, err = strconv.Atoi(parts[0])
	if err != nil {
		println("(Client) Error: Convert to int failed:", err.Error())
		os.Exit(1)
	}

	conn.Close()
}

func handleDownloadFile(tcpServer *net.TCPAddr) bool {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("(Client) Error: Dial failed:", err.Error())
		os.Exit(1)
	}

	fmt.Print("(Client) Enter file name: ")
	var fileName string
	fmt.Scanln(&fileName)
	_, err = conn.Write([]byte("Download " + strconv.Itoa(UserID) + " " + fileName + " "))
	if err != nil {
		println("(Client) Error: Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data

	fileData, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("(Client) Error: Read data failed:", err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(".\\Downloads\\" + fileName , fileData, 0644)
	if err != nil {
		fmt.Println("(Client) Error: Write file failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("(Client) File downloaded successfully")
	fmt.Println("(Client) Do you want to download another file? (y/n)")
	var response string
	fmt.Scanln(&response)
	if response == "y" {
		conn.Close()
		return false
	}

	conn.Close()
	return true
}
