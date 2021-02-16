package websub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kryptn/websub-to-slack/internal/pkg/store"
)

type WebSubConfig struct {
	Name         string
	TopicURL     string
	HubURL       string
	BaseURL      string
	VerifyToken  string
	LeaseSeconds int

	healthcheckUrl string

	mux   *http.ServeMux
	store store.Store
}

func (c *WebSubConfig) params() url.Values {
	params := url.Values{}
	params.Add("hub.mode", "subscribe")
	params.Add("hub.topic", c.TopicURL)
	params.Add("hub.callback", fmt.Sprintf("%s%s", c.BaseURL, c.endpoint()))
	params.Add("hub.verify", "sync")
	params.Add("hub.verify_token", c.VerifyToken)
	params.Add("hub.lease_seconds", strconv.Itoa(c.LeaseSeconds))

	return params
}

func NewWebsubListener(name, topicURL, hubURL, publicHostname, verifyToken string, mux *http.ServeMux, store store.Store, leaseSeconds int, healthcheckUrl string) *WebSubConfig {
	wsc := WebSubConfig{
		Name:        name,
		TopicURL:    topicURL,
		HubURL:      hubURL,
		BaseURL:     publicHostname,
		VerifyToken: verifyToken,

		LeaseSeconds: leaseSeconds,

		mux:   mux,
		store: store,
	}

	return &wsc
}

func (c *WebSubConfig) endpoint() string {
	return fmt.Sprintf("/webhooks/%s", c.Name)
}

func (c *WebSubConfig) renewSubscription() {

	resp, err := http.PostForm(c.HubURL, c.params())
	if err != nil {
		panic(err)
	}

	log.Printf("%s", resp.Status)

	defer resp.Body.Close()
}

func healthy(url string) bool {
	resp, err := http.Get(fmt.Sprintf("%s/healthz", url))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("health check -- %s", resp.Status)
	return err == nil && resp.StatusCode == http.StatusOK

}

func waitForPublic(url string, timeout, attempts int) error {

	timeoutDuration := time.Second * time.Duration(timeout)

	usedAttempts := 0

	for true {

		log.Printf("checking public endpoint")
		if healthy(url) {
			log.Printf("checking public endpoint -- ok")
			return nil
		}

		usedAttempts++
		if usedAttempts >= attempts {
			break
		}

		log.Printf("checking public endpoint -- failed, waiting")
		time.Sleep(timeoutDuration)

	}
	msg := fmt.Sprintf("could not hit public healthcheck url after %d attempts", usedAttempts)
	return errors.New(msg)
}

func (c *WebSubConfig) handleResubscription(ctx context.Context) {
	log.Printf("check public url %s", c.BaseURL)
	if err := waitForPublic(c.BaseURL, 10, 10); err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context was cancelled in sub watcher %s", c.Name)
			return
		case <-ticker.C:
			hasLease, _ := c.store.KeyExists(c.Name)
			if !hasLease {
				log.Printf("renewing %s subscription", c.Name)
				c.renewSubscription()
			}
		}
	}
}

func (c *WebSubConfig) Start(ctx context.Context) <-chan io.Reader {

	message := make(chan io.Reader, 5)

	c.mux.HandleFunc(c.endpoint(), c.MakeHandler(message))

	go c.handleResubscription(ctx)

	return message
}
