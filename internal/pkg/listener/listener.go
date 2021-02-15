package listener

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kryptn/websub-to-slack/internal/pkg/listener/webhook"

	"github.com/kryptn/websub-to-slack/internal/pkg/config"
	"github.com/kryptn/websub-to-slack/internal/pkg/listener/websub"
	"github.com/kryptn/websub-to-slack/internal/pkg/store"
)

type Listener interface {
	Start(ctx context.Context) <-chan io.Reader
}

func ListenersFromConfig(config *config.Config, mux *http.ServeMux, store store.Store) (map[string]Listener, error) {

	listeners := make(map[string]Listener)

	for name, listener := range config.Listeners {
		switch listener.Handler {
		case "websub":
			defaultSeconds := 60 * 60 * 18 
			if listener.LeaseSeconds != 0 {
				defaultSeconds = listener.LeaseSeconds
			}

			healthcheckURL := fmt.Sprintf("%s/healthz", config.PublicHostname)
			listeners[name] = websub.NewWebsubListener(name, listener.TopicURL, listener.HubURL, config.PublicHostname, config.VerifyToken, mux, store, defaultSeconds, healthcheckURL)
		case "webhook":
			listeners[name] = webhook.NewWebhookListener(mux)
		default:
			log.Fatalf("could not identify listener %s", listener.Handler)
		}
	}

	return listeners, nil
}
