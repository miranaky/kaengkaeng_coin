package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miranaky/kaengkaengcoin/blockchain"
	"github.com/miranaky/kaengkaengcoin/utils"
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

type balanceResponse struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
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
			URL:         url("/status"),
			Method:      "GET",
			Discription: "See status of the blockchain",
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
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Discription: "See All Balance by Address",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks())
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
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

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BlockChain())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BlockChain().BalanceByAddress(address)
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		err := json.NewEncoder(rw).Encode(blockchain.BlockChain().TxOutsByAddress(address))
		utils.HandleErr(err)
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation)
	router.HandleFunc("/status", status)
	router.HandleFunc("/block", blocks)
	router.HandleFunc("/block/{hash:[a-f0-9]+}", block)
	router.HandleFunc("/balance/{address}", balance)
	fmt.Printf("Listening REST API on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
