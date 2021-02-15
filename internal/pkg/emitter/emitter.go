package emitter

import (
	"io"

	"github.com/kryptn/websub-to-slack/internal/pkg/emitter/slack"
	"github.com/kryptn/websub-to-slack/internal/pkg/store"

	"github.com/kryptn/websub-to-slack/internal/pkg/config"
)

func EmittersFromConfig(config *config.Config, store store.Store) (map[string]io.Writer, error) {

	emitters := make(map[string]io.Writer)

	for name, emitterConfig := range config.Emitters {

		switch emitterConfig.Handler {
		case "slack":
			slackConfig := slack.NewSlackEmitter(name, emitterConfig.IncomingWebhook, store)
			emitters[name] = slackConfig
		default:
			continue

		}

	}

	return emitters, nil
}
