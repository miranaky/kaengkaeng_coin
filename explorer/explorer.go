package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/miranaky/kaengkaengcoin/blockchain"
)

const (
	port         string = ":4000"
	templatesDir        = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}

}

func Start() {
	templates = template.Must(template.ParseGlob(templatesDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templatesDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}