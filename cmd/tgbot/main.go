package main

import (
	"context"
	"io"
	"os"

	tg "gopkg.in/telebot.v3"

	"github.com/danblok/cleanbg/internal/grpcapi"
	"github.com/danblok/cleanbg/internal/log"
)

func main() {
	log := log.New("local")
	log = log.With("app", "tg_bot")
	log.Warn("logger enabled")

	client, err := grpcapi.Connect("[::]:42069")
	if err != nil {
		log.Error("error", err)
		return
	}

	opts := tg.Settings{
		Token: os.Getenv("TG_TOKEN"),
	}
	bot, err := tg.NewBot(opts)
	if err != nil {
		log.Error("error", err)
		return
	}

	bot.Handle(tg.OnPhoto, func(ctx tg.Context) error {
		rc, err := bot.File(&tg.File{FileID: ctx.Message().Photo.FileID})
		if err != nil {
			return err
		}
		defer rc.Close()

		imageToProcess, err := io.ReadAll(rc)
		if err != nil {
			return err
		}

		img, err := client.Clean(context.Background(), imageToProcess)
		if err != nil {
			return err
		}

		file, err := os.CreateTemp("", "*.jpg")
		if err != nil {
			return err
		}
		defer file.Close()
		defer func() {
			err := os.Remove(file.Name())
			if err != nil {
				log.Error("OnPhoto", "error", err)
			}
		}()

		_, err = file.Write(img)
		if err != nil {
			return err
		}

		res := &tg.Document{File: tg.FromDisk(file.Name()), FileName: file.Name()}
		return ctx.Send(res)
	})

	bot.Start()
}
