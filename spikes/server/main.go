package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.Path, req.Proto)
	for k, values := range req.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
	fmt.Print("\n")
	switch req.Method {
	case "GET":
		for k, values := range req.URL.Query() {
			for _, v := range values {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	case "POST":
		req.ParseForm()
		for k, values := range req.PostForm {
			for _, v := range values {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	}
	fmt.Print("\n\n")
	w.WriteHeader(http.StatusTeapot)
}
