package configfilereader

import (
	"encoding/json"
	"os"
)

/*Configuration struct to hold the Configuration from Config.json*/
type Configuration struct {
	PORT              int    `json:"PORT" bson:"PORT"`
	HOST              string `json:"HOST" bson:"HOST"`
	MongoDBConnString string `json:"MongoDBConnString" bson:"MongoDBConnString"`
}

/*ConfigFilePath full path of Config.json*/
var ConfigFilePath = "./Config/Config.json"

/*ReadConfig Reads Config from Config.json
  Returns Defaults if error is found in reading the Config.json file.*/
func ReadConfig(inConf *Configuration) error {
	// setting defaults
	inConf.PORT = 8080
	inConf.HOST = "localhost"
	inConf.MongoDBConnString = "mongodb://localhost:27017"
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(inConf)
	if err != nil {
		return err
	}
	return err
}
