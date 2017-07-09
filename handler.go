package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type WebhookHandler struct {
	SlackClient *slack.Client
	ChannelID   string
}

func (h WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var event github.PullRequestEvent
	if err := json.Unmarshal(buf, &event); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", string(buf))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if *event.Action == "closed" && *event.PullRequest.Merged {
		phrase := "おめでとうございます、お兄様！お兄様はまたしても不可能を可能にされました！"
		params := slack.NewPostMessageParameters()
		if _, _, err := h.SlackClient.PostMessage(h.ChannelID, phrase, params); err != nil {
			log.Printf("[ERROR] Failed to post message: %s", err)
			return
		}
	}
}
