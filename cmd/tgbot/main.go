package main

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"os"

	tg "gopkg.in/telebot.v3"

	"github.com/danblok/cleanbg/internal/grpcapi"
	"github.com/danblok/cleanbg/internal/log"
)

func main() {
	log := log.New("local")
	log.Warn("logger enabled")
	log = log.WithGroup("tg_bot")

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

		file, err := os.CreateTemp("", "*.png")
		if err != nil {
			return err
		}
		_, err = file.ReadFrom(rc)
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

		processedImage, _, err := image.Decode(bytes.NewReader(img))
		if err != nil {
			return err
		}

		err = png.Encode(file, processedImage)
		if err != nil {
			return err
		}

		res := &tg.Document{File: tg.FromDisk(file.Name()), FileName: file.Name()}
		return ctx.Send(res)
	})

	bot.Start()
}
