package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/miranaky/kaengkaengcoin/blockchain"
)

const (
	// port         string = ":4000"
	templatesDir = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		blockchain.BlockChain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}

}

func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templatesDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templatesDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening at http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
