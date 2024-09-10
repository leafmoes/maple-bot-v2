package router

import (
	"embed"
	"html/template"
	"maple-bot/internal/bot"
	"maple-bot/internal/config"
	"net/http"
)

//go:embed index.html
var indexHtml embed.FS

type HttpHandleFunc func(w http.ResponseWriter, r *http.Request)

var indexTmpl = template.Must(template.ParseFS(indexHtml, "index.html"))

func RegisterRouter(mux *http.ServeMux) {
	// router /
	mux.HandleFunc("/", index(config.App.Domain))
	// router /telegram/
	{
		mux.HandleFunc("/telegram/validate", validate(bot.Token))
		mux.HandleFunc("/telegram/set_webhook", setwebhook())
		mux.HandleFunc(bot.WebhookSubPath, bot.GetUpdater().GetHandlerFunc(bot.WebhookSubPath))
	}
}
