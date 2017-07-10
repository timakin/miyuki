package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/unrolled/render"
)

var ren = render.New()

type WebhookHandler struct {
	SlackClient *slack.Client
	ChannelID   string
}

func (h WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Github-Event") == "ping" {
		ren.JSON(w, http.StatusOK, "")
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		ren.JSON(w, http.StatusMethodNotAllowed, map[string]string{"message": fmt.Sprintf("[ERROR] Invalid method: %s", r.Method)})
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("[ERROR] Failed to read request body: %s", err)})
		return
	}

	var event github.PullRequestEvent
	if err := json.Unmarshal(buf, &event); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", string(buf))
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("[ERROR] Failed to decode json message from slack: %s", string(buf))})
		return
	}

	if event.Action == nil || event.PullRequest == nil {
		ren.JSON(w, http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("[ERROR] Unmarshalized entity was nil")})
		return
	}

	if *event.Action == "closed" && *event.PullRequest.Merged && *event.PullRequest.Base.Ref == "master" {
		phrase := "おめでとうございます、お兄様！お兄様はまたしても不可能を可能にされました！"
		params := slack.NewPostMessageParameters()
		params.Username = "miyuki"
		params.IconURL = "https://github.com/timakin/miyuki/raw/master/miyuki_face.png"
		if _, _, err := h.SlackClient.PostMessage(h.ChannelID, phrase, params); err != nil {
			log.Printf("[ERROR] Failed to post message: %s", err)
			ren.JSON(w, http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("[ERROR] Failed to post message: %s", err)})
			return
		}
	} else {
		log.Printf("[ERROR] Failed to parse a right payload: (action: %s, merged: %s, Ref: %s)", *event.Action, strconv.FormatBool(*event.PullRequest.Merged), *event.PullRequest.Base.Ref)
		ren.JSON(w, http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("[ERROR] Failed to parse a right payload: (action: %s, merged: %s, Ref: %s)", *event.Action, strconv.FormatBool(*event.PullRequest.Merged), *event.PullRequest.Base.Ref)})
		return
	}
}
