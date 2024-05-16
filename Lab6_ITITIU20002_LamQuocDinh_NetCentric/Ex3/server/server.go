package main

import (
	"fmt"
	"net/http"
	"os"
)

func serveIndex(w http.ResponseWriter, req *http.Request) {
	_, err := os.Stat("index.html")
	if os.IsNotExist(err) {
		http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
	http.ServeFile(w, req, "index.html")
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<b>Welcome  to Very Simple Web Server</b> <br>")
	fmt.Fprintf(w, "<i>Welcome  to Very Simple Web Server</i> ")
}
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":9999", nil)
}
