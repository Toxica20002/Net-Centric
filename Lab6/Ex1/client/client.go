package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:9999/hello", nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return
	}

	fmt.Println("Request Message: ")
	fmt.Println(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		return
	}
	// Print the response message
	fmt.Println("Response Message:")
	fmt.Println(resp)
	// Read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}
	// Print the response body
	fmt.Println("Response Body:")
	fmt.Println(string(body))

}
