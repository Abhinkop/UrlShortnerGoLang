package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"UrlShortnerGoLang/internal_pkg/configfilereader"
	"UrlShortnerGoLang/internal_pkg/dbconnect"
	"UrlShortnerGoLang/internal_pkg/randomstringgen"

	"github.com/gorilla/mux"
)

var config configfilereader.Configuration

func main() {
	// Reading the Configuration from Config.json file.
	err := configfilereader.ReadConfig(&config)
	if err != nil {
		log.Println("Error reading config file:", err)
		log.Println("Using defaults ", config)
	}

	var router = mux.NewRouter()
	router.HandleFunc("/api/shortenURL", shortenURLapi).Methods("POST")
	router.HandleFunc("/images/favicon.ico", faviconHandler)
	router.HandleFunc("/shortenURL", shortenURLHandler)

	//We first check for a endpoint and redirect it if it is present.
	router.HandleFunc("/{shortURL}", shortURLHandler)

	// Handle index Page
	router.HandleFunc("/", indexHandler)

	urlTemplate := "%s:%d"
	url := fmt.Sprintf(urlTemplate, config.HOST, config.PORT)
	log.Println("Server running at: \"", url, "\"")
	log.Fatal(http.ListenAndServe(url, router))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./resources/images/favicon.ico")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./resources/HTML/index.html")
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
	template, err := template.ParseFiles("./resources/HTML/error.html")
	err = template.Execute(*w, e)
	log.Println("Error" + e.Error())
	if err != nil {
		log.Fatal(err)
	}
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	session := dbconnect.Connect(config.MongoDBConnString)
	defer dbconnect.Disconnect(session)
	// Create a entry to insert to db
	var entry dbconnect.LookUpDocument
	entry.FullURL = r.FormValue("URL")
	entry.ShortURLEndPoint = randomstringgen.Genarate(6, session)
	err := dbconnect.InsertLookUpEntry(&entry, session)
	if err != nil {
		displayError(&w, err)
		return
	}

	if err != nil {
		displayError(&w, err)
		return
	}
	template, er := template.ParseFiles("./resources/HTML/inserted.html")
	if er != nil {
		displayError(&w, er)
		return
	}
	entry.ShortURLEndPoint = r.Host + "/" + entry.ShortURLEndPoint
	err = template.Execute(w, entry)
	if err != nil {
		displayError(&w, err)
		return
	}
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	session := dbconnect.Connect(config.MongoDBConnString)
	entry, err := dbconnect.GetLookUpEntry(vars["shortURL"], session)
	dbconnect.Disconnect(session)
	// serve endpoint not found
	if err != nil && err.Error() == "not found" {
		template, err := template.ParseFiles("./resources/HTML/endPointNotFound.html")
		if err != nil {
			displayError(&w, err)
			return
		}
		err = template.Execute(w, nil)
		if err != nil {
			displayError(&w, err)
			return
		}

		return
	}
	http.Redirect(w, r, entry.FullURL, http.StatusMovedPermanently)

}

func shortenURLapi(w http.ResponseWriter, r *http.Request) {
	session := dbconnect.Connect(config.MongoDBConnString)
	defer dbconnect.Disconnect(session)
	defer handleErrorAPI(&w)
	decoder := json.NewDecoder(r.Body)
	var req dbconnect.RequestBody
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}
	var entry dbconnect.LookUpDocument
	entry.FullURL = req.URL
	entry.ShortURLEndPoint = randomstringgen.Genarate(6, session)
	err = dbconnect.InsertLookUpEntry(&entry, session)
	if err != nil {
		panic(err)
	}
	entry.ShortURLEndPoint = r.Host + "/" + entry.ShortURLEndPoint
	json.NewEncoder(w).Encode(entry)
}

func handleErrorAPI(w *http.ResponseWriter) {
	if err := recover(); err != nil {
		http.Error(*w, "ERROR", http.StatusInternalServerError)
	}
}
