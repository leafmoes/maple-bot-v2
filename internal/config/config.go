package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var App AppConfig

type AppConfig struct {
	AppName      string             `json:"app_name"`
	Domain       string             `json:"domain"`
	Host         string             `json:"host"`
	Port         string             `json:"port"`
	TelegramBot  TelegramBotConfig  `json:"bot"`
	Github       GithubAppConfig    `json:"github_app"`
	CloudflareKV CloudflareKVConfig `json:"cloudflare_kv"`
	LogLevel     string             `json:"log_level"`
}

type TelegramBotConfig struct {
	Token              string   `json:"token"`
	RunMode            string   `json:"run_mode"`
	WebhookSubPath     string   `json:"webhook_sub_path"`
	WebhookSecretToken string   `json:"webhook_secret_token"`
	WebAppURL          string   `json:"web_app_url"`
	SuperAdmins        []string `json:"super_admins"`
}

type GithubAppConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackPath string `json:"callback_path"`
}

type CloudflareKVConfig struct {
	ApiToken    string `json:"api_token"`
	AccountId   string `json:"account_id"`
	NamespaceId string `json:"namespace_id"`
}

func InitConfig() error {
	var r io.Reader
	botCfgJson := os.Getenv("MAPLE_BOT_CONFIG")
	if botCfgJson == "" {
		wd, _ := os.Getwd()
		bs, err := os.ReadFile(wd + "/config.json")
		if err != nil {
			return fmt.Errorf(`failed to read config file:%w`, err)
		}
		r = bytes.NewReader(bs)
	} else {
		r = strings.NewReader(botCfgJson)
	}
	if err := json.NewDecoder(r).Decode(&App); err != nil {
		return fmt.Errorf(`failed to decode config file:%w`, err)
	}
	return nil
}
