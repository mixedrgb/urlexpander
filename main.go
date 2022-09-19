package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

const ERROR string = "err.gohtml"
const INDEX string = "index.gohtml"
const RESULT string = "result.gohtml"
const PORT string = ":8888"

type le struct {
	Result string
}

var l = le{
	Result: "asdf",
}

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

func getURL(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	if req.Method != "POST" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		log.Println("asdfjkalsdhfkjasdlhf")
	}
	if len(url) > 0 {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
		resp, client_err := client.Get(url)
		if client_err != nil {
			log.Println("Link not correct.")
			http.Redirect(w, req, "/", http.StatusSeeOther)
		} else {
			l = le{
				Result: resp.Header.Get("Location"),
			}
			log.Printf("\"Location\": %s\n", resp.Header.Get("Location"))
		}
	}
}

func result(w http.ResponseWriter, req *http.Request) {
	getURL(w, req)
	tpl.ExecuteTemplate(w, RESULT, l)
}

func docRoot(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, INDEX, nil)
	err := req.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
}

func err(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, ERROR, nil)
}

func main() {
	http.HandleFunc("/", docRoot)
	http.HandleFunc("/result", result)
	http.HandleFunc("/err", err)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Println("Server starting...")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
