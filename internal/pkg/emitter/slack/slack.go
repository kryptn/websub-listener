package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kryptn/websub-to-slack/internal/pkg/store"
	"github.com/mmcdole/gofeed"
)

type Slack struct {
	name            string
	incomingWebhook string
	store           store.Store
}

func (s *Slack) Write(p []byte) (n int, err error) {

	type payload struct {
		Text string `json:"text"`
	}

	ap := gofeed.NewParser()
	feed, _ := ap.Parse(bytes.NewReader(p))

	cacheKey := fmt.Sprintf("%s/%s", s.name, feed.Items[0].GUID)

	exists, _ := s.store.KeyExists(cacheKey)
	if exists {
		return 0, nil
	}

	s.store.SetKey(cacheKey, feed.Items[0].GUID, time.Duration(3)*time.Hour)

	pt := payload{Text: fmt.Sprintf("%s -- %s", feed.Items[0].Author.Name, feed.Items[0].Link)}

	out, _ := json.Marshal(&pt)
	buf := bytes.NewBuffer(out)

	http.Post(s.incomingWebhook, "application/json", buf)

	return 0, nil
}

func NewSlackEmitter(name, webhookURL string) *Slack {

	s := Slack{
		name: name,
		incomingWebhook: webhookURL,
	}

	return &s

}
