package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type SlackMessage struct {
	Text string `json:"text"`
}

// SendSlackNotification envia uma mensagem para o Slack usando um webhook configurado.
func SendSlackNotification(webhookURL, message string) error {
	payload := SlackMessage{Text: message}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send Slack message: %s", resp.Status)
		return err
	}

	log.Println("Slack notification sent successfully")
	return nil
}

// SlackHandler é um endpoint de teste para envio de mensagens no Slack.
func SlackHandler(w http.ResponseWriter, r *http.Request) {
	webhookURL := r.URL.Query().Get("webhook")
	if webhookURL == "" {
		http.Error(w, "Webhook URL is required", http.StatusBadRequest)
		return
	}

	message := r.URL.Query().Get("message")
	if message == "" {
		message = "Default Slack notification message."
	}

	if err := SendSlackNotification(webhookURL, message); err != nil {
		http.Error(w, "Failed to send Slack notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Slack notification sent successfully"))
}
