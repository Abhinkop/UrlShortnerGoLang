package randomStringGen

import (
	"dbConnect"
	"math/rand"
	"time"

	mgo "gopkg.in/mgo.v2"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Genarate(length int, db *mgo.Session) string {
	for {
		rString := generateRandomString(length)
		res := dbConnect.IsShortURLAlreadyPresent(rString, db)
		if res == false {
			return rString
		}
	}

}

func generateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return "/" + string(b)
}
