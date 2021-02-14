package slack_text

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Slack struct {
	name            string
	incomingWebhook string
}

func (s *Slack) Write(p []byte) (n int, err error) {

	type payload struct {
		Text string `json:"text"`
	}

	pt := payload{Text: fmt.Sprintf(string(p))}

	out, _ := json.Marshal(&pt)
	buf := bytes.NewBuffer(out)

	resp, err := http.Post(s.incomingWebhook, "application/json", buf)
	if err != nil {
		log.Printf("uhhh %v", err)
		return 0, err
	}
	log.Printf("webhook %s -- %s", s.name, resp.Status)

	return len(p), nil
}

func NewSlackEmitter(name, webhookURL string) *Slack {

	s := Slack{
		name:            name,
		incomingWebhook: webhookURL,
	}

	return &s

}
