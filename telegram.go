package main

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"io"
)

type File struct {
	name string
	file io.Reader
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *File) Name() string {
	return f.name
}

type Telegram struct {
	bot *telego.Bot
}

func NewTelegram(api, token string) (*Telegram, error) {
	var options []telego.BotOption
	if api != "" {
		options = append(options, telego.WithAPIServer(api))
	}

	bot, err := telego.NewBot(token, options...)
	if err != nil {
		return nil, err
	}

	return &Telegram{bot: bot}, nil
}

func (tg *Telegram) SendMessage(msg string, to int64, markdown, silent bool) error {
	var parseMode string
	if markdown {
		parseMode = telego.ModeMarkdownV2
	}

	_, err := tg.bot.SendMessage(&telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: to},
		Text:      msg,
		ParseMode: parseMode,
	})
	return err
}

func (tg *Telegram) SendDocument(file io.Reader, fileName, caption string, to int64, silent bool) error {
	_, err := tg.bot.SendDocument(&telego.SendDocumentParams{
		ChatID: telego.ChatID{ID: to},
		Document: telego.InputFile{
			File: &File{
				name: fileName,
				file: file,
			},
		},
		Caption:             caption,
		DisableNotification: silent,
	})
	return err
}

func (tg *Telegram) SendPhoto(image io.Reader, caption string, to int64, silent bool) error {
	_, err := tg.bot.SendPhoto(&telego.SendPhotoParams{
		ChatID: telego.ChatID{ID: to},
		Photo: telego.InputFile{
			File: &File{file: image},
		},
		Caption:             caption,
		DisableNotification: silent,
	})
	return err
}

func (tg *Telegram) SendPhotoWithURL(imageURL, caption string, to int64, silent bool) error {
	_, err := tg.bot.SendPhoto(&telego.SendPhotoParams{
		ChatID:              telego.ChatID{ID: to},
		Photo:               telegoutil.FileFromURL(imageURL),
		Caption:             caption,
		DisableNotification: silent,
	})
	return err
}
