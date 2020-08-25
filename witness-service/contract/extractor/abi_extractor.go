package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var jsonFilePath string
	flag.StringVar(&jsonFilePath, "json", "", "input json file")
	flag.Parse()
	if jsonFilePath == "" {
		log.Fatalln("missing input json file path")
	}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	var jsonContent map[string]interface{}
	json.Unmarshal([]byte(byteValue), &jsonContent)
	abiContent, err := json.Marshal(jsonContent["abi"])
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(abiContent))
}
