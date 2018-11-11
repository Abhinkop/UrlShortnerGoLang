package dbConnect

import (
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"log"
)

// The DB Name
const DbName string = "URLShortner"

// The Collection Name
const CollectionName string = "LookUp"

type LookUpDocument struct {
	// Keys in the Collection
	FullURL          string `json:"FullURL" bson:"FullURL"`
	ShortURLEndPoint string `json:"ShortURLEndPoint" bson:"ShortURLEndPoint"`
}

// Connects to a monog db running at uri
func Connect(uri string) *mgo.Session {
	db, err := mgo.Dial(uri)
	if err != nil {
		log.Println("cannot dial mongo", err)
		return nil
	}

	return db
}

// Disconnects to a mongo db pointed by db
func Disconnect(db *mgo.Session) {
	db.Close()
}

// Inserts a entry to the Collection
func InsertLookUpEntry(entry *LookUpDocument, db *mgo.Session) error {
	err := db.DB(DbName).C(CollectionName).Insert(entry)
	if err != nil {
		log.Println("cannot insert into  Collection", err)
	}
	return err
}

// Gets a entry from the collection based on the shortURL
// this only gets one therfore care must be taken that no two same Short URL's must be inserted to  the mongo Db
// Since this is not for  production i have ommited this part.
func GetLookUpEntry(shortURL string, db *mgo.Session) (LookUpDocument, error) {
	var entry LookUpDocument
	err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).One(&entry)
	if err != nil {
		log.Println(err)
	}
	return entry, err
}

// Checks if the  given short URL is already prsent in the collection.
// This is used to keep the short url unique.
func IsShortURLAlreadyPresent(shortURL string, db *mgo.Session) bool {
	shortURLCount, err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).Count()
	if err != nil {
		log.Println(err)
	}
	return shortURLCount != 0
}
