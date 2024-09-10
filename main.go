package main

import (
	"log"
	"maple-bot/internal/bot"
	"maple-bot/internal/config"
	"maple-bot/router"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func main() {
	rm := bot.RunMode
	if rm == "polling" {
		b := bot.GetBot()
		updater := bot.GetUpdater()
		err := updater.StartPolling(b, &ext.PollingOpts{
			DropPendingUpdates: true,
			GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
				Timeout: 9,
				RequestOpts: &gotgbot.RequestOpts{
					Timeout: time.Second * 10,
				},
			},
		})
		if err != nil {
			panic("failed to start polling:" + err.Error())
		}
		log.Printf("%s has been started...\n", b.Username)
		updater.Idle()
	} else if rm == "webhook" {
		b := bot.GetBot()
		updater := bot.GetUpdater()
		err := updater.SetAllBotWebhooks(config.App.Domain+bot.WebhookSubPath, &gotgbot.SetWebhookOpts{
			MaxConnections:     100,
			DropPendingUpdates: true,
			SecretToken:        bot.WebhookSecretToken,
		})
		if err != nil {
			panic("Failed to set bot webhooks: " + err.Error())
		}
		mux := http.NewServeMux()
		router.RegisterRouter(mux)
		server := http.Server{
			Handler: mux,
			Addr:    ":3333",
		}
		log.Printf("%s has been started...\n", b.Username)
		if err := server.ListenAndServe(); err != nil {
			panic("failed to listen and serve: " + err.Error())
		}
	} else {
		panic("Unkown bot run mode: " + rm)
	}
}
