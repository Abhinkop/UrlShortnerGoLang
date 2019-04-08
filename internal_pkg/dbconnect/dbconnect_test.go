package dbconnect

import (
	"testing"

	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Test for dbConnect package.

var mongoAddr = "mongodb://root:123@172.17.0.2"

func TestInsertLookUpEntry(t *testing.T) {
	db := Connect(mongoAddr)
	var c LookUpDocument
	c.FullURL = "https:/www.github.com"
	c.ShortURLEndPoint = "/gh"
	err1 := InsertLookUpEntry(&c, db)

	var entry LookUpDocument
	err := db.DB(DbName).C(CollectionName).Find(bson.M{"ShortURLEndPoint": c.ShortURLEndPoint}).One(&entry)
	if err != nil {
		fmt.Println(err)
	}

	if c.FullURL != entry.FullURL || c.ShortURLEndPoint != entry.ShortURLEndPoint {
		t.Errorf("Inserted: %v, got: %v, want: %v.\n %v", c, entry, c, err1)
	}
	err = db.DB(DbName).C(CollectionName).Remove(bson.M{"ShortURLEndPoint": c.ShortURLEndPoint})
	if err != nil {
		fmt.Println("CleanUp Failed", err)
	}
	Disconnect(db)
}

func TestGetLookUpEntry(t *testing.T) {
	db := Connect(mongoAddr)
	var c LookUpDocument
	c.FullURL = "https:/www.github.com"
	c.ShortURLEndPoint = "/gh"
	InsertLookUpEntry(&c, db)

	d, err := GetLookUpEntry("/gh", db)

	if c.FullURL != d.FullURL || c.ShortURLEndPoint != d.ShortURLEndPoint {
		t.Errorf("Inserted: %v, got: %v, want: %v.\n %v", c, d, c, err)
	}
	err = db.DB(DbName).C(CollectionName).Remove(bson.M{"ShortURLEndPoint": c.ShortURLEndPoint})
	if err != nil {
		fmt.Println("CleanUp Failed", err)
	}
	Disconnect(db)
}

func TestWrongGetLookUpEntry(t *testing.T) {
	db := Connect(mongoAddr)
	_, err := GetLookUpEntry("/gh", db)
	if err.Error() != "not found" {
		t.Errorf(err.Error())
	}
	Disconnect(db)
}

func TestIsShortURLAlreadyPresent(t *testing.T) {
	db := Connect(mongoAddr)
	var c LookUpDocument
	c.FullURL = "https:/www.github.com"
	c.ShortURLEndPoint = "/gh"
	err := InsertLookUpEntry(&c, db)
	res := IsShortURLAlreadyPresent("/gh", db)
	if res != true {
		t.Errorf("IsShortURLAlreadyPresent returned false expecting true")
	}

	err = db.DB(DbName).C(CollectionName).Remove(bson.M{"ShortURLEndPoint": c.ShortURLEndPoint})
	if err != nil {
		fmt.Println("CleanUp Failed", err)
	}

	Disconnect(db)
}

func TestIsShortURLUnique(t *testing.T) {
	db := Connect(mongoAddr)
	res := IsShortURLAlreadyPresent("/hgfkjdh", db)
	if res != false {
		t.Errorf("IsShortURLAlreadyPresent retuned false expected true")
	}
	Disconnect(db)
}
