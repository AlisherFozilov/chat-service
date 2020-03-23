package main

import (
	"github.com/AlisherFozilov/chat-service/pkg/services/fileconnector"
	"io/ioutil"
	"log"
)

//WARNING: this simple client program not pretend to be good or at list
//be normal. Just a program for testing.

// Test file service
func main() {
	fileData, err := ioutil.ReadFile("D:/go.png")
	if err != nil {
		log.Fatal(err)
	}

	connector := fileconnector.NewFileSvcConnector("http://localhost:9999/api/files")
	urls, err := connector.SaveOnFileServiceAndGetUrls(fileData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(urls)
}
