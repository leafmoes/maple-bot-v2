package bot

import (
	"log"
	"maple-bot/internal/bot/handler"
	"maple-bot/internal/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var (
	Token              string
	RunMode            string
	WebhookSubPath     string
	WebhookSecretToken string
	bot                *gotgbot.Bot
	updater            *ext.Updater
)

func init() {
	config.InitConfig()
	Token = config.App.TelegramBot.Token
	RunMode = config.App.TelegramBot.RunMode
	WebhookSubPath = config.App.TelegramBot.WebhookSubPath
	WebhookSecretToken = config.App.TelegramBot.WebhookSecretToken

	var err error
	bot, err = gotgbot.NewBot(Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}
	dispather := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Panicln("an error occurred while handing update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater = ext.NewUpdater(dispather, nil)
	if RunMode == "webhook" {
		err := updater.AddWebhook(bot, Token, &ext.AddWebhookOpts{SecretToken: WebhookSecretToken})
		if err != nil {
			panic("Failed to add bot webhooks to updater: " + err.Error())
		}
	}
	// Register Hanlder /
	handler.RegisterHandlers(dispather)

}

func GetUpdater() *ext.Updater {
	return updater
}

func GetBot() *gotgbot.Bot {
	return bot
}

// type client struct {
// 	rwMux    sync.RWMutex
// 	userData map[int64]map[string]any
// }

// func (c *client) getUserData(ctx *ext.Context, key string) (any, bool) {
// 	c.rwMux.RLock()
// 	defer c.rwMux.RUnlock()

// 	if c.userData == nil {
// 		return nil, false
// 	}

// 	userData, ok := c.userData[ctx.EffectiveUser.Id]
// 	if !ok {
// 		return nil, false
// 	}

// 	v, ok := userData[key]
// 	return v, ok
// }

// func (c *client) setUserData(ctx *ext.Context, key string, val any) {
// 	c.rwMux.Lock()
// 	defer c.rwMux.Unlock()

// 	if c.userData == nil {
// 		c.userData = map[int64]map[string]any{}
// 	}

// 	_, ok := c.userData[ctx.EffectiveUser.Id]
// 	if !ok {
// 		c.userData[ctx.EffectiveUser.Id] = map[string]any{}
// 	}
// 	c.userData[ctx.EffectiveUser.Id][key] = val
// }
