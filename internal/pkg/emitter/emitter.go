package emitter

import (
	"io"

	"github.com/kryptn/websub-to-slack/internal/pkg/emitter/slack"
	"github.com/kryptn/websub-to-slack/internal/pkg/emitter/slack_text"

	"github.com/kryptn/websub-to-slack/internal/pkg/config"
)

func EmittersFromConfig(config *config.Config) (map[string]io.Writer, error) {

	emitters := make(map[string]io.Writer)

	for name, emitterConfig := range config.Emitters {

		switch emitterConfig.Handler {
		case "slack":
			slackConfig := slack.NewSlackEmitter(name, emitterConfig.IncomingWebhook)
			emitters[name] = slackConfig
		case "slack_text":
			emitters[name] = slack_text.NewSlackEmitter(name, emitterConfig.IncomingWebhook)
		default:
			continue

		}

	}

	return emitters, nil
}
