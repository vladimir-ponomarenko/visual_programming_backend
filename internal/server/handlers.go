package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"visprogbackend/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	s := r.Context().Value("server").(*Server)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		http.Error(w, "WebSocket upgrade error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			break
		}

		if messageType == websocket.TextMessage {
			var data models.Message
			if err := json.Unmarshal(message, &data); err != nil {
				log.Println("JSON parsing error:", err)
				if err := conn.WriteMessage(websocket.TextMessage, []byte("JSON parsing error")); err != nil {
					log.Println("Error sending error message:", err)
				}
				continue
			}

			if data.Wcdma == nil {
				data.Wcdma = []models.WcdmaData{}
			}
			if data.Gsm == nil {
				data.Gsm = []models.GsmData{}
			}
			if data.Nr == nil {
				data.Nr = []models.NRData{}
			}
			if data.Lte == nil {
				data.Lte = []models.LteData{}
			}

			fmt.Println("Received data:", data)

			_, err = s.dbpool.Exec(context.Background(), `
				INSERT INTO messages (time, latitude, longitude, altitude, operator, wcdma, gsm, lte, nr)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				data.Time, data.Latitude, data.Longitude, data.Altitude, data.Operator,
				data.Wcdma, data.Gsm, data.Lte, data.Nr)

			if err != nil {
				log.Println("Ошибка записи в базу данных:", err)
				log.Println(err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte("Data received and saved to database")); err != nil {
				log.Println("Error sending message:", err)
				break
			}
		}
	}
}
