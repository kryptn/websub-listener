package webhook

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type WebhookConfig struct {
	mux *http.ServeMux
}

func NewWebhookListener(mux *http.ServeMux) *WebhookConfig {
	return &WebhookConfig{
		mux: mux,
	}
}

func (c *WebhookConfig) Start(ctx context.Context) <-chan io.Reader {
	message := make(chan io.Reader, 5)
	c.mux.HandleFunc("/text", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("error reading body")
		}

		message <- bytes.NewReader(body)
	})

	return message

}
