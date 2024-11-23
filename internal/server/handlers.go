package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"visprogbackend/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	dataMutex sync.Mutex
	messages  []models.Message
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка обновления до WebSocket:", err)
		http.Error(w, "Ошибка обновления до WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			break
		}

		if messageType == websocket.TextMessage {
			var сообщение models.Message
			if err := json.Unmarshal(message, &сообщение); err != nil {
				log.Println("Ошибка парсинга JSON:", err)
				if err := conn.WriteMessage(websocket.TextMessage, []byte("Ошибка парсинга JSON")); err != nil {
					log.Println("Ошибка отправки сообщения об ошибке:", err)
				}
				continue
			}

			fmt.Println("Полученные данные:", сообщение)

			dataMutex.Lock()
			messages = append(messages, сообщение)
			dataMutex.Unlock()

			if err := saveMessagesToFile(); err != nil {
				log.Println("Ошибка записи в файл:", err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte("Данные получены")); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				break
			}
		}
	}
}

func HandleDataRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	dataMutex.Lock()
	defer dataMutex.Unlock()

	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Println("Ошибка маршалинга JSON:", err)
		http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func saveMessagesToFile() error {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(messages)
}
