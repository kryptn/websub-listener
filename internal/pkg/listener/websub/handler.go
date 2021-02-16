package websub

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func (c *WebSubConfig) getHandler(w http.ResponseWriter, r *http.Request) {
	if mode, ok := r.URL.Query()["hub.mode"]; !ok || len(mode) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if mode[0] != "subscribe" && mode[0] != "unsubscribe" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if secret, ok := r.URL.Query()["hub.verify_token"]; !ok || len(secret) == 0 || secret[0] != c.VerifyToken {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("secret didn't match")
		return
	}
	if challenge, ok := r.URL.Query()["hub.challenge"]; ok && len(challenge) > 0 {
		leases, ok := r.URL.Query()["hub.lease_seconds"]
		if !ok {
			leases = []string{"3600"}
		}
		lease, err := strconv.Atoi(leases[0])
		if err != nil {
			lease = 300
		}

		preExipiryLease := lease - (lease / 20)

		err = c.store.SetKey(c.Name, lease, time.Duration(preExipiryLease)*time.Second)
		if err != nil {
			fmt.Println(err)
		}
		//s.Cache.SetLease(s.Slug, lease)
		log.Printf("setting lease %s for %d, originally %d", c.Name, preExipiryLease, lease)

		log.Printf("got challenge: %s -- responding", challenge[0])
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, challenge[0])
		return
	}
}

func (c *WebSubConfig) postHandler(event chan<- io.Reader) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("oh yeah we're writing")

		w.WriteHeader(http.StatusOK)

		ap := gofeed.NewParser()
		feed, _ := ap.Parse(r.Body)

		if len(feed.Items) == 0 {
			return
		}

		cacheKey := fmt.Sprintf("%s/%s", c.Name, feed.Items[0].GUID)

		if exists, _ := c.store.KeyExists(cacheKey); !exists {
			c.store.SetKey(cacheKey, feed.Items[0].GUID, time.Duration(12)*time.Hour)

			log.Printf("setting %s for %s -- %s", c.Name, feed.Items[0].GUID, cacheKey)
			type payload struct {
				Text string `json:"text"`
			}

			message := fmt.Sprintf("%s -- %s", feed.Items[0].Author.Name, feed.Items[0].Link)
			event <- strings.NewReader(message)
		}
	}
}

func (c *WebSubConfig) MakeHandler(event chan<- io.Reader) http.HandlerFunc {

	log.Printf("created %s handler", c.Name)

	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			c.getHandler(w, r)
		case http.MethodPost:
			c.postHandler(event)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}
