package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pebbe/zmq4"
)

// Notifier define o contrato para envio de log a destinos externos.
type Notifier interface {
	Notify(entry *LogEntry)
}

// ExternalNotifier envia o log para uma URL externa (via HTTP)
// e opcionalmente para um endpoint ZMQ.
type ExternalNotifier struct {
	externalURL string
	zmqSocket   *zmq4.Socket
}

// NewExternalNotifier cria uma nova instância de ExternalNotifier.
func NewExternalNotifier(url string, zmqEndpoint string) *ExternalNotifier {
	var socket *zmq4.Socket
	if zmqEndpoint != "" {
		var err error
		socket, err = zmq4.NewSocket(zmq4.PUSH)
		if err != nil {
			log.Printf("Erro ao criar socket ZMQ: %v", err)
		} else {
			socket.Connect(zmqEndpoint)
		}
	}
	return &ExternalNotifier{
		externalURL: url,
		zmqSocket:   socket,
	}
}

// Notify envia o log via HTTP (se externalURL estiver configurada)
// e via ZMQ (se o socket estiver ativo).
func (n *ExternalNotifier) Notify(entry *LogEntry) {
	if n.externalURL != "" {
		data, _ := json.Marshal(entry)
		_, err := http.Post(n.externalURL, "application/json", strings.NewReader(string(data)))
		if err != nil {
			log.Printf("Erro ao enviar log para %s: %v", n.externalURL, err)
		}
	}
	if n.zmqSocket != nil {
		data, _ := json.Marshal(entry)
		n.zmqSocket.Send(string(data), 0)
	}
}

// DiscordNotifier envia logs para um webhook do Discord.
type DiscordNotifier struct {
	webhook string
}

// NewDiscordNotifier cria uma nova instância de DiscordNotifier.
func NewDiscordNotifier(webhook string) *DiscordNotifier {
	return &DiscordNotifier{
		webhook: webhook,
	}
}

// Notify formata a mensagem e a envia via HTTP POST para o Discord.
func (n *DiscordNotifier) Notify(entry *LogEntry) {
	if n.webhook == "" {
		return
	}
	message := fmt.Sprintf("**[%s] %s**\n%s",
		entry.Level,
		entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		entry.Message)
	payload := map[string]string{"content": message}
	jsonPayload, _ := json.Marshal(payload)
	_, err := http.Post(n.webhook, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		log.Printf("Erro ao enviar log para Discord: %v", err)
	}
}
