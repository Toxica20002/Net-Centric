/**
 * Author: Yin Lin
 * Server side of the guessing game
**/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

var UserList = make(map[int]int)

const NUMBER_OF_USER = 10000

type Users struct {
	Users []User `json:"users"`
}

type gameRoom struct {
	timeStart      time.Time
	numberOfPlayer int
	gameRoomID     int
	startGame      bool
	gameRoomWord   string
	guessWord      string
	playerPoint    map[int]int
	playerTurn     map[int]int
	turn           int
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Usernames string `json:"username"`
	Passwords string `json:"password"`
}

var gameRoomList = make(map[int]gameRoom)

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
				rand.Seed(time.Now().UnixNano())
				userID := rand.Intn(NUMBER_OF_USER) + 1
				_, exist := UserList[userID]
				for exist {
					userID = rand.Intn(NUMBER_OF_USER) + 1
					_, exist = UserList[userID]
				}
				response := strconv.Itoa(userID) + " "
				UserList[userID] = 0
				conn.Write([]byte(response))
				validUser = true
				break
			}
		}

		if !validUser {
			conn.Write([]byte("Invalid "))
		}

	} else if parts[0] == "logout" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}
		delete(UserList, id)
		conn.Write([]byte("(Server) Logout successfully "))
	} else if parts[0] == "waiting" {
		//Check available game room
		//If no available game room, create a new game room
		//If there is available game room, join the game room
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}

		if UserList[id] != 0 {
			userGameRoomID := UserList[id]
			userGameRoom := gameRoomList[userGameRoomID]
			timeNow := time.Now()
			duration := timeNow.Sub(userGameRoom.timeStart)
			if duration.Seconds() >= 10 {
				if userGameRoom.numberOfPlayer == 1 {
					UserList[id] = 0
					delete(gameRoomList, userGameRoomID)
					response := strconv.Itoa(userGameRoom.numberOfPlayer) + " 0 1 "
					conn.Write([]byte(response))
				} else {
					userGameRoom.startGame = true
					gameRoomList[userGameRoomID] = userGameRoom
					response := strconv.Itoa(userGameRoom.numberOfPlayer) + " 1 0 "
					conn.Write([]byte(response))
				}
			} else {
				response := strconv.Itoa(userGameRoom.numberOfPlayer) + " 0 0 "
				conn.Write([]byte(response))
			}

			return
		}

		//Check available game room
		availableGameRoom := false
		for _, value := range gameRoomList {
			if !value.startGame {
				availableGameRoom = true
				UserList[id] = value.gameRoomID
				break
			}
		}

		if !availableGameRoom {
			rand.Seed(time.Now().UnixNano())
			gameRoomID := rand.Intn(NUMBER_OF_USER) + 1
			_, exist := gameRoomList[gameRoomID]
			for exist {
				gameRoomID = rand.Intn(NUMBER_OF_USER) + 1
				_, exist = gameRoomList[gameRoomID]
			}
			UserList[id] = gameRoomID
			randomWord := randomdata.Adjective()
			len := len(randomWord)
			guessWord := ""
			for i := 0; i < len; i++ {
				guessWord += "_"
			}
			fmt.Println("Random word: " + randomWord)
			gameRoomList[gameRoomID] = gameRoom{time.Now(), 1, gameRoomID, false, randomWord, guessWord, make(map[int]int), make(map[int]int), 1}
			gameRoomList[gameRoomID].playerPoint[id] = 0
			gameRoomList[gameRoomID].playerTurn[id] = 1
			response := "1 0 0 "
			conn.Write([]byte(response))
		} else {
			for key, value := range gameRoomList {
				if !value.startGame {
					value.numberOfPlayer++
					gameRoomList[key] = value
					UserList[id] = value.gameRoomID
					response := strconv.Itoa(value.numberOfPlayer) + " 0 0 "
					conn.Write([]byte(response))
					value.playerPoint[id] = 0
					value.playerTurn[id] = value.numberOfPlayer
					break
				}
			}

		}

	} else if parts[0] == "start" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}

		userGameRoomID := UserList[id]
		response := strconv.Itoa(userGameRoomID) + " "
		conn.Write([]byte(response))
	} else if parts[0] == "checkturn" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}

		userGameRoomID := UserList[id]
		userGameRoom := gameRoomList[userGameRoomID]

		if userGameRoom.guessWord == userGameRoom.gameRoomWord {
			response := "end "
			conn.Write([]byte(response))
		}

		//fmt.Println(userGameRoom.turn)
		if userGameRoom.playerTurn[id] == userGameRoom.turn {
			conn.Write([]byte("yes "))
		} else {
			conn.Write([]byte("no "))
		}
	} else if parts[0] == "nextturn" {
		//fmt.Println("Next turn")
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}
		userGameRoomID := UserList[id]
		userGameRoom := gameRoomList[userGameRoomID]
		userGameRoom.turn++
		if userGameRoom.turn > userGameRoom.numberOfPlayer {
			userGameRoom.turn = 1
		}
		gameRoomList[userGameRoomID] = userGameRoom
	} else if parts[0] == "getpoint" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}
		userGameRoomID := UserList[id]
		userGameRoom := gameRoomList[userGameRoomID]
		response := strconv.Itoa(userGameRoom.playerPoint[id]) + " "
		conn.Write([]byte(response))
	} else if parts[0] == "getword" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}
		userGameRoomID := UserList[id]
		userGameRoom := gameRoomList[userGameRoomID]
		response := userGameRoom.guessWord + " "
		conn.Write([]byte(response))
	} else if parts[0] == "guessword" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}

		guessChar := string(parts[2])
		guessChar = strings.TrimSpace(guessChar)
		//fmt.Println("Guess char: " + guessChar)
		userGameRoomID := UserList[id]
		userGameRoom := gameRoomList[userGameRoomID]
		guessWord := userGameRoom.guessWord
		flag := false
		for i := 0; i < len(userGameRoom.gameRoomWord); i++ {
			if string(userGameRoom.gameRoomWord[i]) == guessChar && string(guessWord[i]) == "_" {
				guessWord = guessWord[:i] + guessChar + guessWord[i+1:]
				userGameRoom.playerPoint[id] += 10
				flag = true
			}
		}

		userGameRoom.guessWord = guessWord
		gameRoomList[userGameRoomID] = userGameRoom

		if !flag {
			response := "0 "
			conn.Write([]byte(response))
		} else {
			response := "1 "
			conn.Write([]byte(response))
		}

	} else if parts[0] == "end" {
		userID := string(parts[1])
		userID = strings.TrimSpace(userID)
		id, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println("(Server): " + string(err.Error()))
			os.Exit(1)
		}
		userGameRoomID := UserList[id]
		userGameRoom := gameRoomList[userGameRoomID]
		UserList[id] = 0
		// find the winner
		winner := 0
		maxPoint := 0
		count := 0
		for key, value := range userGameRoom.playerPoint {
			if value > maxPoint {
				maxPoint = value
				winner = key
				count = 1
			} else if value == maxPoint {
				count++
			}
		}
		answerWord := userGameRoom.gameRoomWord
		response := strconv.Itoa(winner) + " " + strconv.Itoa(maxPoint) + " " + strconv.Itoa(count) + " " + answerWord + " "
		conn.Write([]byte(response))

	}

	// close conn
	conn.Close()
}
