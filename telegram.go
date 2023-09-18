package main

import (
	"io"
	"time"

	tb "gopkg.in/telebot.v3"
)

type Telegram struct {
	bot *tb.Bot
}

func NewTelegram(api, token string) (*Telegram, error) {
	bot, err := tb.NewBot(tb.Settings{
		URL:    api,
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	return &Telegram{bot: bot}, nil
}

func (tg *Telegram) SendMessage(msg string, to int64, markdown, silent bool) error {
	opt := &tb.SendOptions{DisableNotification: silent}
	if markdown {
		opt.ParseMode = tb.ModeMarkdownV2
	}

	_, err := tg.bot.Send(tb.ChatID(to), msg, opt)
	return err
}

func (tg *Telegram) SendFile(file io.Reader, fileName, mime, caption string, to int64, silent bool) error {
	_, err := tg.bot.Send(tb.ChatID(to), &tb.Document{
		File:     tb.File{FileReader: file},
		Caption:  caption,
		MIME:     mime,
		FileName: fileName,
	}, &tb.SendOptions{DisableNotification: silent})
	return err
}

func (tg *Telegram) SendImage(image io.Reader, caption string, to int64, silent bool) error {
	_, err := tg.bot.Send(tb.ChatID(to), &tb.Photo{
		File:    tb.File{FileReader: image},
		Caption: caption,
	}, &tb.SendOptions{DisableNotification: silent})
	return err
}
