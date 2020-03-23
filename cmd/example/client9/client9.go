package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlisherFozilov/chat-service/pkg/services/messaging"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

//WARNING: this simple client program not pretend to be good or at list
//be normal. Just a program for testing.

// silly client
func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8888/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	selfID, receiverID := int64(1), int64(2)
	fmt.Println("Input your ID")
	_, err = fmt.Scan(&selfID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Input your friend's (receiver) ID")
	_, err = fmt.Scan(&receiverID)
	if err != nil {
		log.Fatal(err)
	}

	idJSON := messaging.IdJSON{Id: selfID}
	idData, err := json.Marshal(idJSON)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, idData)
	if err != nil {
		log.Fatal(err)
	}

	msg := messaging.CommonMessage{
		Type:       messaging.Text,
		ReceiverID: receiverID,
	}

	go write(&msg, conn)
	read(conn)
}

func read(conn *websocket.Conn) {
	msg := messaging.CommonMessage{}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg.Data)
	}
}

func write(msg *messaging.CommonMessage, conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		msg.Data = scanner.Text()
		data, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}

		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
