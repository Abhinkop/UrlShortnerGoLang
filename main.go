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
	template, err := template.ParseFiles("./htmlTemplates/endPointNotFound.html")
	if err != nil {
		displayError(&w, err)
		return
	}
	err = template.Execute(w, nil)
	if err != nil {
		displayError(&w, err)
		return
	}
}

func displayError(w *http.ResponseWriter, e error) {
	template, err := template.ParseFiles("./htmlTemplates/error.html")
	err = template.Execute(*w, e)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./htmlTemplates/index.html")
	if err != nil {
		displayError(&w, err)
		return
	}
	err = template.Execute(w, nil)
	if err != nil {
		displayError(&w, err)
		return
	}
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	session := dbConnect.Connect("localhost")
	// Create a entry to insert to db
	var entry dbConnect.LookUpDocument
	entry.FullURL = r.FormValue("URL")
	entry.ShortURLEndPoint = randomStringGen.Genarate(6, session)
	err := dbConnect.InsertLookUpEntry(&entry, session)
	if err != nil {
		displayError(&w, err)
		return
	}
	dbConnect.Disconnect(session)

	if err != nil {
		displayError(&w, err)
		return
	}
	template, er := template.ParseFiles("./htmlTemplates/inserted.html")
	if er != nil {
		displayError(&w, er)
		return
	}
	entry.ShortURLEndPoint = r.Host + entry.ShortURLEndPoint
	err = template.Execute(w, entry)
	if err != nil {
		displayError(&w, err)
		return
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./images/favicon.ico")
}

func main() {
	fmt.Println("Starting the server on ", Host)
	mux := http.NewServeMux()
	httpHandler := NewHttpRedirectHandler(mux)
	http.ListenAndServe(":8080", httpHandler)
}

func NewHttpRedirectHandler(fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			// server root path
			{
				indexHandler(w, r)
				return
			}
		case "/shortenURL":
			// serve shorten URL API
			{
				shortenURLHandler(w, r)
				return
			}
		case "/images/favicon.ico":
			// Server favicon
			{
				faviconHandler(w, r)
				return
			}
		default:
			{
				session := dbConnect.Connect("localhost")
				entry, err := dbConnect.GetLookUpEntry(r.URL.Path, session)
				dbConnect.Disconnect(session)
				// serve endpoint not found
				if err != nil && err.Error() == "not found" {
					redirect(w, r)
					return
				} else {
					http.Redirect(w, r, entry.FullURL, http.StatusMovedPermanently)
				}

			}

		}
	}
}
