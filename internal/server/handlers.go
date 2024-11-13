package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"visprogbackend/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // TODO: В продакшене необходимо настроить CheckOrigin для безопасности!  Например: return r.Header.Get("Origin") == "http://localhost:3000"
}

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

			файл, err := os.OpenFile("данные.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("Ошибка открытия файла:", err)
				continue
			}
			defer файл.Close()

			jsonДанные, err := json.Marshal(сообщение)
			if err != nil {
				log.Println("Ошибка маршалинга JSON:", err)
				continue
			}

			if _, err := файл.Write(jsonДанные); err != nil {
				log.Println("Ошибка записи в файл:", err)
				continue
			}

			if _, err := файл.Write([]byte("\n")); err != nil {
				log.Println("Ошибка записи новой строки в файл:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte("Данные получены")); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				break
			}
		}
	}
}
