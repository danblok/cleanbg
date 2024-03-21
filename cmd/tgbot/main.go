package main

import (
	"os"

	tg "gopkg.in/telebot.v3"

	"github.com/danblok/cleanbg/internal/api"
	"github.com/danblok/cleanbg/internal/bot"
	"github.com/danblok/cleanbg/internal/log"
)

func main() {
	log := log.New("local")
	log.Warn("logger enabled")

	client, err := api.Connect(os.Getenv("GRPC_ADDR"))
	if err != nil {
		log.Error("grpc failed", "error", err)
		return
	}

	opts := tg.Settings{
		Token: os.Getenv("TG_TOKEN"),
	}
	bot, err := bot.New(client, opts)
	if err != nil {
		log.Error("bot failed", "error", err)
	}
	bot.Start()
}
