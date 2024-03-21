package bot

import (
	"context"
	"io"
	"os"

	tg "gopkg.in/telebot.v3"

	"github.com/danblok/cleanbg/internal/types"
)

type Bot struct {
	bot    *tg.Bot
	client types.Cleaner
}

func New(client types.Cleaner, opts tg.Settings) (*Bot, error) {
	bot, err := tg.NewBot(opts)
	if err != nil {
		return nil, err
	}

	b := &Bot{bot: bot, client: client}
	bot.Handle("/start", b.handleGreeting)
	bot.Handle(tg.OnPhoto, b.handleClear)

	return b, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) handleGreeting(ctx tg.Context) error {
	return ctx.Send(
		"Hello, it's your best bot to remove background from your dirty images.\n\nUsing me is stupidly easy. All what you need to do is to send me a picture and I will return the same picture without the background. By the way, various formats are supported. Enjoy :3",
	)
}

func (b *Bot) handleClear(ctx tg.Context) error {
	rc, err := b.bot.File(&tg.File{FileID: ctx.Message().Photo.FileID})
	if err != nil {
		return err
	}
	defer rc.Close()

	imageToProcess, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	img, err := b.client.Clean(context.Background(), imageToProcess)
	if err != nil {
		return err
	}

	file, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		return err
	}
	defer file.Close()
	defer func() {
		os.Remove(file.Name())
	}()

	_, err = file.Write(img)
	if err != nil {
		return err
	}

	res := &tg.Document{File: tg.FromDisk(file.Name()), FileName: file.Name()}
	return ctx.Send(res)
}
