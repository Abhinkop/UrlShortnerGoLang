package dbConnect

import (
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"log"
)

const DbName string = "URLShortner"
const CollectionName string = "LookUp"

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
	err := db.DB(DbName).C(CollectionName).Insert(entry)
	if err != nil {
		log.Println("cannot insert into  Collection", err)
	}
	return err
}

func GetLookUpEntry(shortURL string, db *mgo.Session) (LookUpDocument, error) {
	var entry LookUpDocument
	err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).One(&entry)
	if err != nil {
		log.Println(err)
	}
	return entry, err
}

func IsShortURLAlreadyPresent(shortURL string, db *mgo.Session) bool {
	shortURLCount, err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).Count()
	if err != nil {
		log.Println(err)
	}
	return shortURLCount != 0
}
