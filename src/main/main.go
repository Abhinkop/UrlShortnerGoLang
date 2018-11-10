package main

import (
	"dbConnect"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"randomStringGen"
)

var Host = "localhost:8080"

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Endpoint not found")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html")
	err = template.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	session := dbConnect.Connect("localhost")
	var entry dbConnect.LookUpDocument
	entry.FullURL = r.FormValue("URL")
	entry.ShortURLEndPoint = randomStringGen.Genarate(6, session)
	err := dbConnect.InsertLookUpEntry(&entry, session)
	dbConnect.Disconnect(session)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	template, er := template.ParseFiles("./inserted.html")
	if er != nil {
		fmt.Println(er)
		os.Exit(-1)
	}
	entry.ShortURLEndPoint = r.Host + entry.ShortURLEndPoint
	err = template.Execute(w, entry)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../images/favicon.ico")
}

func main() {
	fmt.Println("Starting the server on ", Host)
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/shortenURL", shortenURLHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)
	xxx := NewHttpRedirectHandler(mux)
	http.ListenAndServe(":8080", xxx)
}

func NewHttpRedirectHandler(fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session := dbConnect.Connect("localhost")
		entry, err := dbConnect.GetLookUpEntry(r.URL.Path, session)
		dbConnect.Disconnect(session)
		if err != nil {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, entry.FullURL, http.StatusMovedPermanently)
		}

	}
}
