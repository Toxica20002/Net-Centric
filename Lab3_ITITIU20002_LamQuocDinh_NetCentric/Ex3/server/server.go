/**
 * Author: Yin Lin
 * Server side of the guessing game
**/

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

var NumberOfUser = 0
var UserList = make(map[int]int)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID  int    `json:"id"`
	Name string `json:"name"`
	Usernames string `json:"username"`
	Passwords string `json:"password"`
}


func main() {

	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("(Server): " + string(err.Error()))
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
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
		fmt.Println("(Server): " + string(err.Error()))
		os.Exit(1)
	}

	request := string(buffer[:])
	parts := strings.Split(request, " ")

	if parts[0] == "login" {
		username := string(parts[1])
		password := string(parts[2])

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		jsonFile, err := os.Open("users.json")
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var users Users
		json.Unmarshal(byteValue, &users)
		
		validUser := false

		for _, user := range users.Users {
			if user.Usernames == username && user.Passwords == password {
				NumberOfUser++
				response := strconv.Itoa(NumberOfUser) + " "
                UserList[NumberOfUser] = user.ID
				conn.Write([]byte(response))
				validUser = true
				break
			}
		}

		if !validUser {
			conn.Write([]byte("Invalid "))
		}
		
	} else if parts[0] == "Download" {
		userID, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}

		fileName := parts[2]

        if _, ok := UserList[userID]; ok {
            defer conn.Close()

            fileData, err := ioutil.ReadFile(".\\Stuff\\" + fileName)
            if err != nil {
                fmt.Println("(Server): " + string(err.Error()))
                os.Exit(1)
            }
            conn.Write(fileData)
        } else {
            conn.Write([]byte("Invalid "))
        }

	} else if(parts[0] == "List") {
        response := ""
        files, err := ioutil.ReadDir(".\\Stuff")
        if err != nil {
            fmt.Println("(Server): " + string(err.Error()))
            os.Exit(1)
        }

        for _, file := range files {
            response += file.Name() + " "
        }

        conn.Write([]byte(response))
    }


	// close conn
	conn.Close()
}
