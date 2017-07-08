package config

import (
	"encoding/json"
	"log"
	"os"
)

//Configuration represents every configuration made in config.json
type Configuration struct {
	LogDir  string
	LogFile string
}

//LoadConfig takes the filename of the json config file, parses it and returns a pointer to a Configuration struct
func LoadConfig(configFile string) *Configuration {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	decodingErr := decoder.Decode(&configuration)
	if decodingErr != nil {
		log.Fatalln(decodingErr)
	}
	log.Printf("logDir: %s\n", configuration.LogDir)
	return &configuration
}
