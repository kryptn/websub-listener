package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

func (s *Subscription) SendSubscribe() {
	params := url.Values{}
	params.Add("hub.mode", "subscribe")
	params.Add("hub.topic", s.TopicURL)
	params.Add("hub.callback", fmt.Sprintf("%s%s", viper.GetString("publicUrl"), s.endpoint()))
	params.Add("hub.verify", "sync")
	params.Add("hub.verify_token", s.VerifyToken)

	resp, err := http.PostForm(viper.GetString("websubSubscribeHost"), params)
	if err != nil {
		panic(err)
	}

	log.Printf("%s", resp.Status)

	defer resp.Body.Close()
}
