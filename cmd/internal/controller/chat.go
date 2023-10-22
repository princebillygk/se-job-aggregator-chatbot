package controller

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/princebillygk/se-job-aggregator-chatbot/cmd/internal/utility"
)

var msngrVerfToken string

func init() {
	var ok bool
	msngrVerfToken, ok = os.LookupEnv("MESSENGER_WEBHOOK_INTEGRATION_TOKEN")
	if !ok {
		panic("Messenger Verification Webhook doesn't exists")
	}
}

// Messenger is a controller for messaging services
type Messenger struct {
}

func (c Messenger) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleNotification(w, r)
	case "GET":
		handleVerification(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body["object"] != "page" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	utility.DampVar(body)
	w.WriteHeader(200)
}

func handleVerification(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	mode, token, challenge := query.Get("hub.mode"), query.Get("hub.verify_token"), query.Get("hub.challenge")

	if mode != "subscribe" && token != msngrVerfToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(challenge))
}