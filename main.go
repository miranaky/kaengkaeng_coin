package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":4000"

type URL string

type URLDiscription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Discription string `json:"discription"`
	Payload     string `json:"payload,omitempty"`
}

func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDiscription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Discription: "See Documentation",
		},
		{
			URL:         URL("/block"),
			Method:      "GET",
			Discription: "See All Blocks",
		},
		{
			URL:         URL("/block"),
			Method:      "POST",
			Discription: "Add A Block",
			Payload:     "data:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func main() {
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
