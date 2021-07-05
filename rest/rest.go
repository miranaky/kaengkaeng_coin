package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miranaky/kaengkaengcoin/blockchain"
)

//
var port string

type url string

type urlDiscription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Discription string `json:"discription"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDiscription{
		{
			URL:         url("/"),
			Method:      "GET",
			Discription: "See Documentation",
		},
		{
			URL:         url("/block"),
			Method:      "GET",
			Discription: "See All Blocks",
		},
		{
			URL:         url("/block"),
			Method:      "POST",
			Discription: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/block/{hash}"),
			Method:      "GET",
			Discription: "See A Block",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// json.NewEncoder(rw).Encode(blockchain.BlockChain().AllBlocks())
		return
	case "POST":
		// var addBlockBody addBlockBody
		// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		// blockchain.BlockChain().AddBlock(addBlockBody.Message)
		// rw.WriteHeader(http.StatusCreated)
		return
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation)
	router.HandleFunc("/block", blocks)
	router.HandleFunc("/block/{hash:[a-f0-9]+}", block)
	fmt.Printf("Listening REST API on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
