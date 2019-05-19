package dbconnect

import (
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"log"
)

/*DbName The DB Name */
const DbName string = "URLShortner"

/*CollectionName The Collection Name */
const CollectionName string = "LookUp"

/*LookUpDocument the strcture of the DBObject */
type LookUpDocument struct {
	// Keys in the Collection
	FullURL          string `json:"FullURL" bson:"FullURL"`
	ShortURLEndPoint string `json:"ShortURLEndPoint" bson:"ShortURLEndPoint"`
}

/*RequestBody struct to hold the Post request Body*/
type RequestBody struct {
	URL string `json:"URL" bson:"URL"`
}

/*Connect Connects to a monog db running at uri */
func Connect(uri string) *mgo.Session {
	if uri == "" {
		log.Fatalln("Error empty Connection String to mongo Db servers")
	}
	db, err := mgo.Dial(uri)
	if err != nil {
		log.Fatalln("cannot dial mongo:", uri, "\n with error:", err)
	}

	return db
}

/*Disconnect Disconnects to a mongo db pointed by db */
func Disconnect(db *mgo.Session) {
	db.Close()
}

/*InsertLookUpEntry Inserts a entry to the Collection */
func InsertLookUpEntry(entry *LookUpDocument, db *mgo.Session) error {
	err := db.DB(DbName).C(CollectionName).Insert(entry)
	if err != nil {
		log.Println("cannot insert into  Collection", err)
	}
	return err
}

/*GetLookUpEntry Gets a entry from the collection based on the shortURL
this only gets one therfore care must be taken that no two same Short URL's must be inserted to  the mongo Db
Since this is not for  production i have ommited this part. */
func GetLookUpEntry(shortURL string, db *mgo.Session) (LookUpDocument, error) {
	var entry LookUpDocument
	err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).One(&entry)
	return entry, err
}

/*IsShortURLAlreadyPresent Checks if the  given short URL is already prsent in the collection.
This is used to keep the short url unique. */
func IsShortURLAlreadyPresent(shortURL string, db *mgo.Session) bool {
	shortURLCount, err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": shortURL}).Count()
	if err != nil {
		log.Println(err)
	}
	return shortURLCount != 0
}
