package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/BMaksim/TestTaskGolang/app/api"
)

var configPath string = "../config.json"

func main() {
	config := api.NewConfig()
	str, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(str, &config)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(config)
	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}

}
