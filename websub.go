package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

func (s Subscription) params() url.Values {
	params := url.Values{}
	params.Add("hub.mode", "subscribe")
	params.Add("hub.topic", s.TopicURL)
	params.Add("hub.callback", fmt.Sprintf("%s%s", viper.GetString("publicUrl"), s.endpoint()))
	params.Add("hub.verify", "sync")
	params.Add("hub.verify_token", s.VerifyToken)

	return params
}

func (s *Subscription) RenewSubscription() {

	resp, err := http.PostForm(viper.GetString("websubSubscribeHost"), s.params())
	if err != nil {
		panic(err)
	}

	log.Printf("%s", resp.Status)

	defer resp.Body.Close()
}

func healthy(url string) bool {
	resp, err := http.Get(fmt.Sprintf("%s/healthz", url))
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

func (c *Config) WatchSubs(timeout int) {
	log.Printf("check public url")
	if err := waitForPublic(c.PublicURL, 10, 10); err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	waitTime := time.Second * time.Duration(timeout)

	go func() {
		for true {
			now := time.Now().Unix()
			for name, listener := range c.Listeners {
				lease, ok := c.Cache.Leases[name]
				if !ok || now > lease {
					log.Printf("renewing %s subscription", name)
					listener.RenewSubscription()
				} else {
					log.Printf("%s subscription has %d seconds left", name, lease-now)
				}
			}

			time.Sleep(waitTime)
		}
	}()

}
