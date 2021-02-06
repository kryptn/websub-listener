package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
)

func (s Subscription) getHandler(w http.ResponseWriter, r *http.Request) {
	if mode, ok := r.URL.Query()["hub.mode"]; !ok || len(mode) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if mode[0] != "subscribe" && mode[0] != "unsubscribe" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if secret, ok := r.URL.Query()["hub.verify_token"]; !ok || len(secret) == 0 || secret[0] != s.VerifyToken {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("secret didn't match")
		return
	}
	if challenge, ok := r.URL.Query()["hub.challenge"]; ok && len(challenge) > 0 {
		leases, ok := r.URL.Query()["hub.lease_seconds"]
		if !ok {
			leases = []string{"3600"}
		}
		lease := leases[0]
		// spew.Dump(s)
		s.Cache.SetLease(s.Slug, lease)
		log.Printf("setting lease %s for %s", s.Slug, lease)

		log.Printf("got challenge: %s -- responding", challenge[0])
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, challenge[0])
		return
	}
}

func (s Subscription) postHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("oh yeah we're writing")

	w.WriteHeader(http.StatusOK)

	ap := gofeed.NewParser()
	feed, _ := ap.Parse(r.Body)

	if s.Cache.ShouldAct(s.Slug, feed.Items[0].GUID) {
		log.Printf("setting %s for %s", s.Slug, feed.Items[0].GUID)
		type payload struct {
			Text string `json:"text"`
		}

		p := payload{Text: fmt.Sprintf("%s -- %s", feed.Items[0].Author.Name, feed.Items[0].Link)}

		out, _ := json.Marshal(&p)
		buf := bytes.NewBuffer(out)

		http.Post(s.PostURL, "application/json", buf)
	}
}

func (s Subscription) MakeHandler() http.HandlerFunc {

	log.Printf("created %s handler", s.Slug)

	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			s.getHandler(w, r)
		case http.MethodPost:
			s.postHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}

func (c *Config) RegisterListeners(mux *http.ServeMux) {
	for name, listener := range c.Listeners {
		mux.HandleFunc(listener.endpoint(), listener.MakeHandler())
		log.Printf("registered %s:%s -- %s", name, listener.Slug, listener.endpoint())
	}

	mux.HandleFunc("/status", c.Cache.CacheStatusHandler)
}
