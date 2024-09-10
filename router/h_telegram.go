package router

import (
	"encoding/json"
	"fmt"
	"maple-bot/internal/bot"
	"maple-bot/internal/config"
	"net/http"
	"net/url"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func setwebhook() HttpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := bot.GetUpdater().SetAllBotWebhooks(config.App.Domain+bot.WebhookSubPath, &gotgbot.SetWebhookOpts{
			MaxConnections:     100,
			DropPendingUpdates: true,
			SecretToken:        bot.WebhookSecretToken,
		})
		if err != nil {
			w.Write([]byte("failed to set bot webhooks:" + err.Error()))
		}
		w.Write([]byte("webhook was set"))
	}
}

func validate(token string) HttpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authQuery, err := url.ParseQuery(r.Header.Get("X-Auth"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("validation failed; failed to parse auth query: " + err.Error()))
		}

		ok, err := ext.ValidateWebAppQuery(authQuery, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("validation failed; error: " + err.Error()))
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("validation failed; data cannot be trusted."))
			return
		}
		var u gotgbot.User
		err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("validation failed; failed to unmarshal user: " + err.Error()))
			return
		}
		w.Write([]byte(fmt.Sprintf("validation success; user '%s' is authenticated (id: %d).", u.FirstName, u.Id)))
	}
}
