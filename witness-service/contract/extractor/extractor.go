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
	var content string
	flag.StringVar(&jsonFilePath, "json", "", "input json file")
	flag.StringVar(&content, "content", "abi", "input json file")
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
	if content != "abi" && content != "bin" {
		log.Fatalln("unknown content", content)
	}
	c, err := json.Marshal(jsonContent[content])
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(c))
}
