package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kryptn/websub-to-slack/internal/pkg/config"
	"github.com/kryptn/websub-to-slack/internal/pkg/emitter"
	"github.com/kryptn/websub-to-slack/internal/pkg/listener"
	"github.com/kryptn/websub-to-slack/internal/pkg/store"
)

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

	log.Print("awaiting signal")
	<-done
	log.Print("exiting")
}

func distribute(ctx context.Context, r <-chan io.Reader, w ...io.Writer) {

	for {

		select {
		case <-ctx.Done():
			return
		case reader := <-r:

			writer := io.MultiWriter(w...)

			if _, err := io.Copy(writer, reader); err != nil {
				log.Fatalf("error on writes %v", err)
			}

		}
	}
}

func main() {

	ctx := context.Background()

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	store := store.StoreFromConfig(config)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	listeners, err := listener.ListenersFromConfig(config, mux, store)
	if err != nil {
		log.Fatal(err)
	}

	emitters, err := emitter.EmittersFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	for listenerName, emitterNames := range config.Wires {
		listener, ok := listeners[listenerName]
		if !ok {
			log.Fatal("could not find listener %s", listenerName)
		}

		var writers []io.Writer

		for _, emitterName := range emitterNames {
			emitter, ok := emitters[emitterName]
			if !ok {
				log.Fatal("could not find emitter %s", emitter)
			}

			writers = append(writers, emitter)
		}

		s := listener.Start(ctx)

		go distribute(ctx, s, writers...)

	}

	go func() {
		log.Fatal(http.ListenAndServe(":8080", mux))
	}()

	awaitSignals()
}
