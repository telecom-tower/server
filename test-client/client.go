package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	log.Println("Ready...")
	http.DefaultClient.Get("http://127.0.0.1:8080/open")
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/update", nil)
	if err != nil {
		log.Fatal(err)
	}
	conn.WriteMessage(websocket.BinaryMessage, []byte{0, 0, 0, 255, 0, 0, 255, 0, 0})
	conn.ReadMessage()

	http.DefaultClient.Get("http://127.0.0.1:8080/close")

}
