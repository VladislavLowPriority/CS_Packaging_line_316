package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

var (
	opcClient *opcua.Client
	mutex     sync.Mutex
	nodeID    = `ns=3;s="Tag_1"`
)

func main() {
	// Подключаемся к OPC UA серверу
	endpoint := "opc.tcp://192.168.0.1:4840"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	opcClient, err = opcua.NewClient(endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}

	if err := opcClient.Connect(ctx); err != nil {
		log.Fatalf("Ошибка подключения к OPC UA серверу: %v", err)
	}
	defer opcClient.Close(ctx)

	// Запускаем веб-сервер
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/toggle", handleToggle)
	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := readValue(ctx)
	if err != nil {
		http.Error(w, "Ошибка чтения", http.StatusInternalServerError)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Управление переменной</title>
    <script>
        function toggleValue() {
            fetch('/toggle')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('value').textContent = data.value;
                })
                .catch(error => console.error('Ошибка:', error));
        }
    </script>
</head>
<body>
    <h1>Текущее значение: <span id="value">{{.}}</span></h1>
    <button onclick="toggleValue()">Переключить</button>
</body>
</html>`

	tmplParsed, _ := template.New("index").Parse(tmpl)
	tmplParsed.Execute(w, value)
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mutex.Lock()
	defer mutex.Unlock()

	value, err := readValue(ctx)
	if err != nil {
		http.Error(w, "Ошибка чтения", http.StatusInternalServerError)
		return
	}

	newValue := !value.(bool)
	if err := writeValue(ctx, newValue); err != nil {
		http.Error(w, "Ошибка записи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"value": %t}`, newValue)
}

func readValue(ctx context.Context) (interface{}, error) {
	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга NodeID: %w", err)
	}

	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{{NodeID: id}},
	}

	res, err := opcClient.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения: %w", err)
	}

	return res.Results[0].Value.Value(), nil
}

func writeValue(ctx context.Context, value bool) error {
	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return fmt.Errorf("ошибка парсинга NodeID: %v", err)
	}

	req := &ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        ua.MustVariant(value),
				},
			},
		},
	}

	res, err := opcClient.Write(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка записи: %v", err)
	}

	if res.Results[0] != ua.StatusOK {
		return fmt.Errorf("ошибка статуса: %s", res.Results[0].Error())
	}

	return nil
}
