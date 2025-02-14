package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/websocket"
    "google.golang.org/protobuf/proto"
    pb "main/chat" // Подключаем gRPC-протофайлы
)

// Обработчик WebSocket-соединений
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Ошибка WebSocket:", err)
        return
    }
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Ошибка чтения:", err)
            break
        }

        var req pb.MessageRequest
        if err := proto.Unmarshal(msg, &req); err != nil {
            log.Println("Ошибка декодирования:", err)
            continue
        }

        log.Printf("Получено сообщение от %s: %s", req.Username, req.Text)

        res := &pb.MessageResponse{Status: "OK "+  req.Username + " " + req.Text}
        resData, _ := proto.Marshal(res)
        conn.WriteMessage(websocket.BinaryMessage, resData)
    }
}

// Раздача `index.html`
func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html")
}

func main() {
    // Проверяем наличие HTML-файла
    if _, err := os.Stat("static/index.html"); os.IsNotExist(err) {
        log.Fatal("Файл `static/index.html` не найден! Создайте его в папке `static/`")
    }

    // Раздаем статические файлы (включая chat.proto)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Обработчики
    http.HandleFunc("/", serveIndex) // Главная страница
    http.HandleFunc("/ws", handleConnection) // WebSocket

    log.Println("Сервер запущен на порту 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
