package dbconnect

import (
	"UrlShortnerGoLang/internal_pkg/configfilereader"
	"strings"
	"testing"

	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Test for dbConnect package.

var conf configfilereader.Configuration

func TestInsertLookUpEntry(t *testing.T) {
	setupServer(t)
	db := Connect(conf.MongoDBConnString)
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
	setupServer(t)
	db := Connect(conf.MongoDBConnString)
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
	setupServer(t)
	db := Connect(conf.MongoDBConnString)
	_, err := GetLookUpEntry("/gh", db)
	if strings.Compare(err.Error(), "not found") != 0 {
		t.Error(err)
	}
	Disconnect(db)
}

func TestIsShortURLAlreadyPresent(t *testing.T) {
	setupServer(t)
	db := Connect(conf.MongoDBConnString)
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
	setupServer(t)
	db := Connect(conf.MongoDBConnString)
	res := IsShortURLAlreadyPresent("/hgfkjdh", db)
	if res != false {
		t.Errorf("IsShortURLAlreadyPresent retuned false expected true")
	}
	Disconnect(db)
}

func setupServer(t *testing.T) {
	configfilereader.ConfigFilePath = "./Config.test.json"
	err := configfilereader.ReadConfig(&conf)
	if err != nil {
		t.Error("Error reading config file Config.test.json")
	}
}
