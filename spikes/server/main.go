package main

import (
	"encoding/xml"
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
	fmt.Println("Host:", req.Host)
	for k, values := range req.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
	fmt.Print("\n")

	var payload map[string][]string
	switch req.Method {
	case "GET":
		payload = req.URL.Query()
	case "POST":
		req.ParseForm()
		payload = req.PostForm
	}

	for k, values := range payload {
		for _, v := range values {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
	fmt.Print("\n\n")

	if action, ok := payload["Action"]; ok && len(action) > 0 {
		if action[0] == "GetQueueURL" {
			if ok := getQueueURL(w, req, payload); ok {
				return
			}
		}
	}

	w.WriteHeader(http.StatusTeapot)
}

type GetQueueURLResponse struct {
	QueueURL string `xml:"GetQueueUrlResponse>GetQueueUrlResult>QueueUrl"`
}

func getQueueURL(w http.ResponseWriter, req *http.Request, payload map[string][]string) bool {
	if queue, ok := payload["QueueName"]; ok && len(queue) > 0 {
		enc := xml.NewEncoder(w)
		enc.Encode(GetQueueURLResponse{
			QueueURL: "http://127.0.0.1:8080/" + queue[0],
		})
		return true
	}
	return false
}
