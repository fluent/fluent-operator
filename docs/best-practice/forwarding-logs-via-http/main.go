package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Log    string `json:"log"`
	Stream string `json:"stream"`
	Time   string `json:"time"`
}

// Sample HTTP receiver for this demo
func main() {
	h := func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		var msg []Message
		err = json.Unmarshal(b, &msg)
		if err != nil {
			log.Printf(err.Error())
			return
		}

		for _, item := range msg {
			log.Printf("log=%s, stream=%s, time=%s\n", item.Log, item.Stream, item.Time)
		}
	}

	http.HandleFunc("/", h)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
