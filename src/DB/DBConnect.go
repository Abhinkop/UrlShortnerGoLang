package main

import (
	//"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"log"
)

type LookUp struct {
	FullURL          string `json:"FullURL" bson:"FullURL"`
	ShortURLEndPoint string `json:"ShortURLEndPoint" bson:"ShortURLEndPoint"`
}

func Connect(uri string) *mgo.Session {
	db, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
		return nil
	}

	return db
}

func Disconnect(db *mgo.Session) {
	db.Close() // clean up when we’re done
}

func Insert(entry *LookUp, dbName, collectionName *string, db *mgo.Session) {
	err := db.DB(*dbName).C(*collectionName).Insert(entry)
	if err != nil {
		fmt.Println("cannot insert into  Collection", err)
	}
}

func main() {
	// connect to the database
	db := Connect("localhost")
	// if err != nil {
	// 	log.Fatal("cannot dial mongo", err)
	// }
	//defer Disconnect(db)
	var c LookUp
	//c.ID = bson.NewObjectID()
	c.FullURL = "https:/www.github.com"
	c.ShortURLEndPoint = "/gh"
	var dbName string = "URLShortner"
	var collectionName string = "LookUp"
	Insert(&c, &dbName, &collectionName, db)
	// err := db.DB("URLShortner").C("LookUp").Insert(&c)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	var d LookUp
	err := db.DB("URLShortner").C("LookUp").Find(bson.M{"ShortURLEndPoint": "/gh"}).One(&d)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(d)
}
