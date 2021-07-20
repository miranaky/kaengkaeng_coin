package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miranaky/kaengkaengcoin/blockchain"
	"github.com/miranaky/kaengkaengcoin/p2p"
	"github.com/miranaky/kaengkaengcoin/utils"
	"github.com/miranaky/kaengkaengcoin/wallet"
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

type myWalletResponse struct {
	Address string `json:"address"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addPeerPayload struct {
	Address, Port string
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
		{
			URL:         url("/ws"),
			Method:      "GET",
			Discription: "Upgrade to WebSockets",
		},
		{
			URL:         url("/peer"),
			Method:      "GET",
			Discription: "connect with upgrade websocket",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		newBlock := blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
		p2p.BroadcastNewBlock(newBlock)
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
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.BlockChain(), rw)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.BlockChain())
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		err := json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.BlockChain()))
		utils.HandleErr(err)
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	// utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
	blockchain.MempoolStatus(blockchain.Mempool(), rw)
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{"not enough fund"})
	}
	rw.WriteHeader(http.StatusCreated)
	p2p.BroadcastNewTx(tx)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})

}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}

}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/block", blocks).Methods("GET", "POST")
	router.HandleFunc("/block/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peer", peers).Methods("GET", "POST")
	fmt.Printf("Listening REST API on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
