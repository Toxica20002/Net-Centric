/**
 * Author: Yin Lin
 * Client side of the guessing game
**/

package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
	TIME_OUT = 30
)

var UserID = 0
var RoomID = 0

type waitingRoom struct {
	numberOfPlayer int
	startGame      int
	runOutOfTime   int
}

type winner struct {
	winnerID int
	winnerPoint int
	count int
	answerWord string
}

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		fmt.Println("(Anonymous Client) Error: ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	handleLogin(tcpServer)

	for {
		fmt.Printf("(Client %d) Do you want to play the game? (yes/no): ", UserID)
		var input string
		fmt.Scanln(&input)
		if input == "yes" {
			waitingRoomResponse := handleWaitingRoom(tcpServer)

			for {
				if waitingRoomResponse.startGame == 1 {
					fmt.Printf("(Client %d) Game is starting\n", UserID)
					handleStartGame(tcpServer)
					fmt.Printf("(Client %d) You are in room %d\n", UserID, RoomID)
					for {
						checkTurn := handleCheckTurn(tcpServer)
						if checkTurn == 1 {
							point := handleGetPoint(tcpServer)
							guessWord := handleGetWord(tcpServer)
							fmt.Printf("(Client %d) It's your turn\n", UserID)
							fmt.Printf("(Client %d) Your point: %d. Current Word: %s\n", UserID, point, guessWord)
							fmt.Printf("(Client %d) Enter your guess: ", UserID)
							guessChar := make(chan rune)
							reader := bufio.NewReader(os.Stdin)
							go func() {
								temp, _, _ := reader.ReadRune()
								guessChar <- temp
							}()

							select {
							case <-time.After(TIME_OUT * time.Second):{
								fmt.Printf("(Client %d) Time out\n", UserID)
							}
							case temp := <-guessChar:{
								result := handleGuessWord(tcpServer, temp)

								if result {
									fmt.Printf("(Client %d) Correct guess\n", UserID)
									point := handleGetPoint(tcpServer)
									guessWord := handleGetWord(tcpServer)
									fmt.Printf("(Client %d) Your point: %d. Current Word: %s\n", UserID, point, guessWord)
								} else {
									fmt.Printf("(Client %d) Character is not exist or has already been guessed \n", UserID)
									handleNextTurn(tcpServer)
								}
							}
							}
							

						} else if checkTurn == 0 {
							point := handleGetPoint(tcpServer)
							guessWord := handleGetWord(tcpServer)
							waitingSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
							waitingSpinner.Prefix = fmt.Sprintf("(Server) You are not in turn. Your point: %d. Current Word: %s", point, guessWord)
							waitingSpinner.Start()
							time.Sleep(1 * time.Second)
							waitingSpinner.Stop()
						} else {
							point := handleGetPoint(tcpServer)
							winner := handleEndGame(tcpServer)

							fmt.Printf("(Client %d) Answer: %s\n", UserID, winner.answerWord)
							if winner.winnerID == UserID && winner.count == 1 {
								fmt.Printf("(Client %d) You are the winner. Your point: %d\n", UserID, winner.winnerPoint)
							} else if winner.winnerPoint == point {
								fmt.Printf("(Client %d) Draw. Your point: %d\n", UserID, point)
							} else {
								fmt.Printf("(Client %d) You lose. Your point: %d\n", UserID, point)
							}
							break
						}
					}
					break
				}
				if waitingRoomResponse.runOutOfTime == 1 {
					fmt.Printf("(Client %d) Game is cancelled\n", UserID)
					break
				}
				waitingSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
				waitingSpinner.Prefix = fmt.Sprintf("(Client %d) Waiting for %d players to join", UserID, waitingRoomResponse.numberOfPlayer)
				waitingSpinner.Start()
				time.Sleep(1 * time.Second)
				waitingSpinner.Stop()
				waitingRoomResponse = handleWaitingRoom(tcpServer)
			}

		} else if input == "no" {
			fmt.Printf("(Client %d) Do you want to log out? (yes/no): ", UserID)
			fmt.Scanln(&input)
			if input == "yes" {
				UserID = 0
				handleLogout(tcpServer)
				break
			}
		} else {
			fmt.Println("Invalid input")
		}
	}
}

func handleLogin(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Anonymous Client) Error: Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Print("(Anonymous Client) Enter your username: ")
	var username string
	fmt.Scanln(&username)
	fmt.Print("(Anonymous Client) Enter your password: ")
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
		fmt.Printf("(Anonymous Client) Error: Write data failed: %s\n", err.Error())
		os.Exit(1)
	}

	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Anonymous Client) Error: Read data failed: %s\n", err.Error())
		os.Exit(1)
	}
	parts := strings.Split(string(received), " ")
	if parts[0] == "Invalid" {
		fmt.Println("(Server) Response: Invalid username or password")
		os.Exit(1)
	}

	UserID, err = strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("(Anonymous Client) Error: Convert to int failed: %s\n", err.Error())
		os.Exit(1)
	}

	conn.Close()
}

func handleLogout(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("logout " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	fmt.Println(response)

	conn.Close()
	UserID = 0
}

func handleWaitingRoom(tcpServer *net.TCPAddr) waitingRoom {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("waiting " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	parts := strings.Split(response, " ")

	conn.Close()

	numberOfPlayer, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to int failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	startGame, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to bool failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	runOutOfTime, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to bool failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	return waitingRoom{numberOfPlayer, startGame, runOutOfTime}
}

func handleStartGame(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("start " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	parts := strings.Split(response, " ")
	RoomID, err = strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to int failed: %s\n", UserID, err.Error())
	}

	conn.Close()
}

func handleCheckTurn(tcpServer *net.TCPAddr) int {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("checkturn " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	parts := strings.Split(response, " ")

	conn.Close()

	if parts[0] == "yes" {
		return 1
	} else if parts[0] == "no" {
		return 0
	} else {
		return 2
	}
}

func handleNextTurn(tcpServer *net.TCPAddr) {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("nextturn " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	conn.Close()
}

func handleGetPoint(tcpServer *net.TCPAddr) int {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("getpoint " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	parts := strings.Split(response, " ")

	point, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to int failed: %s\n", UserID, err.Error())
	}

	conn.Close()
	return point
}

func handleGetWord(tcpServer *net.TCPAddr) string {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("getword " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	parts := strings.Split(response, " ")

	conn.Close()
	return parts[0]
}

func handleGuessWord(tcpServer *net.TCPAddr, guessChar rune) bool {
	//fmt.Println("Guessing Char: ", string(guessChar))
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("guessword " + strconv.Itoa(UserID) + " " + string(guessChar) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	response := string(received[:])
	parts := strings.Split(response, " ")

	conn.Close()

	return parts[0] == "1"
}

func handleEndGame(tcpServer *net.TCPAddr) winner {
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		fmt.Printf("(Client %d) Error: Dial failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("end " + strconv.Itoa(UserID) + " "))
	if err != nil {
		fmt.Printf("(Client %d) Error: Write data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		fmt.Printf("(Client %d) Error: Read data failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	conn.Close()
	response := string(received[:])
	parts := strings.Split(response, " ")
	userWinnerID := parts[0]
	userWinnerID = strings.TrimSpace(userWinnerID)
	id, err := strconv.Atoi(userWinnerID)
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to int failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}
	winnerPoint := parts[1]
	winnerPoint = strings.TrimSpace(winnerPoint)
	point, err := strconv.Atoi(winnerPoint)
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to int failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	count := parts[2]
	count = strings.TrimSpace(count)
	countInt, err := strconv.Atoi(count)
	if err != nil {
		fmt.Printf("(Client %d) Error: Convert to int failed: %s\n", UserID, err.Error())
		os.Exit(1)
	}

	answerWord := parts[3]
	answerWord = strings.TrimSpace(answerWord)

	return winner{id, point, countInt, answerWord}
}
