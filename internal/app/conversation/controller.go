package conversation

import (
	"log"

	"github.com/gorilla/websocket"
)

func HandleAudioStream(conn *websocket.Conn) {
	audioStream := make(chan []byte)

	go func() {
		defer close(audioStream)

		for audioData := range audioStream {
			err := conn.WriteMessage(websocket.BinaryMessage, audioData)
			if err != nil {
				log.Printf("Error writing audio data to websocket: %v", err)
				return
			}
		}
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from websocket: %v", err)
			return
		}

		switch messageType {
		case websocket.BinaryMessage:
			audioStream <- p // TODO process audio data here instead of just passing it through
		case websocket.CloseMessage:
			return
		}
	}
}
