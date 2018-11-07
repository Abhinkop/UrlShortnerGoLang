package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../images/favicon.ico")
}

func main() {
	fmt.Println("Starting the server on ", Host)
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/shortenURL", redirect)
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
