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
	entry, _ := dbConnect.GetLookUpEntry("/fb", session)
	var entry1 dbConnect.LookUpDocument
	entry1.FullURL = r.FormValue("URL")
	entry1.ShortURLEndPoint = "/" + randomStringGen.Genarate(6)
	err := dbConnect.InsertLookUpEntry(&entry1, session)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, entry.FullURL)
	fmt.Fprintln(w, randomStringGen.Genarate(6))
	fmt.Fprintln(w, r.FormValue("URL"))
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
	var xxxx = make(map[string]string)
	xxxx["/fb"] = "https://facebook.com"

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		url1, ok := xxxx[r.URL.Path]
		if ok != true {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url1, http.StatusMovedPermanently)
		}

	}
}
