package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var debug = false

type Config struct {
	PublicURL           string                  `mapstructure:"publicUrl"`
	Listeners           map[string]Subscription `mapstructure:"listener"`
	Destinations        map[string]string       `mapstructure:"destinations"`
	VerifyToken         string                  `mapstructure:"verifyToken"`
	WebsubSubscribeHost string                  `mapstructure:"websubSubscribeHost"`
	Cache               *Cache
}

type Subscription struct {
	Slug        string
	TopicURL    string
	VerifyToken string
	Parser      string
	PostURL     string
	Destination string
	Cache       *Cache
}

func (s *Subscription) endpoint() string {
	return fmt.Sprintf("/webhooks/%s", s.Slug)
}

func awaitSignals() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

func main() {

	log.Printf("oh we're gonna try ")

	config := getConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	config.RegisterListeners(mux)

	go func() {
		log.Fatal(http.ListenAndServe(":9090", mux))
	}()

	for true {
		resp, err := http.Get(fmt.Sprintf("%s/healthz", config.PublicURL))
		if err != nil {
			log.Printf("Waiting until URL is public")
			time.Sleep(time.Second * 5)
		}

		if resp.StatusCode == http.StatusOK {
			break
		}
	}

	time.Sleep(time.Second * 2)

	for name, listener := range config.Listeners {

		log.Printf("sending subscribe for %s", name)
		listener.SendSubscribe()
	}

	awaitSignals()
}
