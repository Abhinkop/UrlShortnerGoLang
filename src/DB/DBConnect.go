package DBConnect

import (
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"log"
)

const dbName string = "URLShortner"
const collectionName string = "LookUp"

type LookUpDocument struct {
	FullURL          string `json:"FullURL" bson:"FullURL"`
	ShortURLEndPoint string `json:"ShortURLEndPoint" bson:"ShortURLEndPoint"`
}

func Connect(uri string) *mgo.Session {
	db, err := mgo.Dial(uri)
	if err != nil {
		log.Println("cannot dial mongo", err)
		return nil
	}

	return db
}

func Disconnect(db *mgo.Session) {
	db.Close()
}

func InsertLookUpEntry(entry *LookUpDocument, db *mgo.Session) error {
	err := db.DB(dbName).C(collectionName).Insert(entry)
	if err != nil {
		log.Println("cannot insert into  Collection", err)
	}
	return err
}

func GetLookUpEntry(shortURL string, db *mgo.Session) (LookUpDocument, error) {
	var entry LookUpDocument
	err := db.DB(dbName).C(collectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).One(&entry)
	if err != nil {
		log.Println(err)
	}
	return entry, err
}
