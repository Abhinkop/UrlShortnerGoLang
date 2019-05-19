package configfilereader

import (
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	ConfigFilePath = "./Config.json"
	var conf Configuration
	err := ReadConfig(&conf)
	if err != nil {
		t.Error("TestReadConfig Failed with error: ", err)
	}
	var expected Configuration
	expected.PORT = 0
	expected.HOST = "TestHOSTAddr"
	expected.MongoDBConnString = "TestConnString"
	if strings.Compare(conf.HOST, expected.HOST) != 0 ||
		conf.PORT != expected.PORT ||
		strings.Compare(conf.MongoDBConnString, expected.MongoDBConnString) != 0 {
		t.Error("ReadConfig test failed. \n Expected:", expected, "\nGot:", conf)
	}
}
