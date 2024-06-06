package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID            int
	Username      string
	Password      string
	FirstName     string
	LastName      string
	Email         string
	Phone         string
	DOB           string
	Avatar        string
	Country       string
	City          string
	StreetName    string
	StreetAddress string
}

type User_API struct {
	ID                    int          `json:"id"`
	UID                   string       `json:"uid"`
	Password              string       `json:"password"`
	FirstName             string       `json:"first_name"`
	LastName              string       `json:"last_name"`
	Username              string       `json:"username"`
	Email                 string       `json:"email"`
	Avatar                string       `json:"avatar"`
	Gender                string       `json:"gender"`
	PhoneNumber           string       `json:"phone_number"`
	SocialInsuranceNumber string       `json:"social_insurance_number"`
	DateOfBirth           string       `json:"date_of_birth"`
	Employment            Employment   `json:"employment"`
	Address               Address      `json:"address"`
	CreditCard            CreditCard   `json:"credit_card"`
	Subscription          Subscription `json:"subscription"`
}

type Employment struct {
	Title    string `json:"title"`
	KeySkill string `json:"key_skill"`
}

type Address struct {
	City          string      `json:"city"`
	StreetName    string      `json:"street_name"`
	StreetAddress string      `json:"street_address"`
	ZipCode       string      `json:"zip_code"`
	State         string      `json:"state"`
	Country       string      `json:"country"`
	Coordinates   Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type CreditCard struct {
	CCNumber string `json:"cc_number"`
}

type Subscription struct {
	Plan          string `json:"plan"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	Term          string `json:"term"`
}

func addUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	resp, errAPI := http.Get("https://random-data-api.com/api/v2/users")
	if errAPI != nil {
		fmt.Println(errAPI)
		return
	}
	defer resp.Body.Close()

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		fmt.Println(errRead)
		return
	}

	var userAPI User_API
	errUnmarshal := json.Unmarshal(body, &userAPI)
	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
		return
	}

	if user.Username == "" {
		user.Username = userAPI.Username
	}

	if user.Password == "" {
		user.Password = userAPI.Password
	}

	if user.FirstName == "" {
		user.FirstName = userAPI.FirstName
	}

	if user.LastName == "" {
		user.LastName = userAPI.LastName
	}

	if user.Email == "" {
		user.Email = userAPI.Email
	}

	if user.Phone == "" {
		user.Phone = userAPI.PhoneNumber
	}

	if user.DOB == "" {
		user.DOB = userAPI.DateOfBirth
	}

	if user.Avatar == "" {
		user.Avatar = userAPI.Avatar
	}

	if user.Country == "" {
		user.Country = userAPI.Address.Country
	}

	if user.City == "" {
		user.City = userAPI.Address.City
	}

	if user.StreetName == "" {
		user.StreetName = userAPI.Address.StreetName
	}

	if user.StreetAddress == "" {
		user.StreetAddress = userAPI.Address.StreetAddress
	}

	_, err := db.Exec(`INSERT INTO users (username, password, first_name, last_name, email, phone, dob, avatar, country, city, street_name, street_address) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Username, user.Password, user.FirstName, user.LastName, user.Email, user.Phone, user.DOB, user.Avatar, user.Country, user.City, user.StreetName, user.StreetAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User added successfully"})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	var user User
	err := db.QueryRow(`SELECT * FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.DOB, &user.Avatar, &user.Country, &user.City, &user.StreetName, &user.StreetAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})

}

func getUserByEmailAddress(c *gin.Context) {
	email := c.Param("email")
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	var user User
	err := db.QueryRow(`SELECT * FROM users WHERE email = ?`, email).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.DOB, &user.Avatar, &user.Country, &user.City, &user.StreetName, &user.StreetAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func updateUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	_, err := db.Exec(`UPDATE users SET username = ?, password = ?, first_name = ?, last_name = ?, email = ?, phone = ?, dob = ?, avatar = ?, country = ?, city = ?, street_name = ?, street_address = ? WHERE id = ?`,
		user.Username, user.Password, user.FirstName, user.LastName, user.Email, user.Phone, user.DOB, user.Avatar, user.Country, user.City, user.StreetName, user.StreetAddress, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User updated successfully"})
}

func deleteUserByID(c *gin.Context) {
	id := c.Param("id")
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User deleted successfully"})
}

func getAllUsers(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.DOB, &user.Avatar, &user.Country, &user.City, &user.StreetName, &user.StreetAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func deleteUserByEmail(c *gin.Context) {
	email := c.Param("email")
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	_, err := db.Exec(`DELETE FROM users WHERE email = ?`, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User deleted successfully"})
}

func getUsernamesByASC(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT username FROM users ORDER BY username ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var usernames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		usernames = append(usernames, username)
	}

	c.JSON(http.StatusOK, gin.H{"usernames": usernames})
}

func getUsernamesByDESC(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT username FROM users ORDER BY username DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var usernames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		usernames = append(usernames, username)
	}

	c.JSON(http.StatusOK, gin.H{"usernames": usernames})
}

func getFirstNameByASC(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT first_name FROM users ORDER BY first_name ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var firstNames []string
	for rows.Next() {
		var firstName string
		err := rows.Scan(&firstName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		firstNames = append(firstNames, firstName)
	}

	c.JSON(http.StatusOK, gin.H{"first_names": firstNames})
}

func getLastNameByASC(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT last_name FROM users ORDER BY last_name ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var lastNames []string
	for rows.Next() {
		var lastName string
		err := rows.Scan(&lastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		lastNames = append(lastNames, lastName)
	}

	c.JSON(http.StatusOK, gin.H{"last_names": lastNames})
}

func getFirstNameByDESC(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT first_name FROM users ORDER BY first_name DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var firstNames []string
	for rows.Next() {
		var firstName string
		err := rows.Scan(&firstName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		firstNames = append(firstNames, firstName)
	}

	c.JSON(http.StatusOK, gin.H{"first_names": firstNames})
}

func getLastNameByDESC(c *gin.Context) {
	db, errDB := sql.Open("mysql", "root:123456789@tcp(localhost:3306)/lab9")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT last_name FROM users ORDER BY last_name DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var lastNames []string
	for rows.Next() {
		var lastName string
		err := rows.Scan(&lastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		lastNames = append(lastNames, lastName)
	}

	c.JSON(http.StatusOK, gin.H{"last_names": lastNames})

}

func main() {
	router := gin.Default()

	//API
	router.POST("/user/add", addUser)
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})
	router.GET("/user/get/id/:id", getUserByID)
	router.GET("/user/get/all", getAllUsers)
	router.GET("/user/get/email/:email", getUserByEmailAddress)
	router.PUT("/user/update/id/:id", updateUserByID)
	router.DELETE("/user/delete/id/:id", deleteUserByID)
	router.DELETE("/user/delete/email/:email", deleteUserByEmail)
	router.GET("/user/get/username/asc", getUsernamesByASC)
	router.GET("/user/get/username/desc", getUsernamesByDESC)
	router.GET("/user/get/first_name/asc", getFirstNameByASC)
	router.GET("/user/get/first_name/desc", getFirstNameByDESC)
	router.GET("/user/get/last_name/asc", getLastNameByASC)
	router.GET("/user/get/last_name/desc", getLastNameByDESC)

	// Run the server
	router.Run(":8080")
}
