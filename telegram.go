package main

import (
	"fmt"
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

	bot.Handle("/start", sendID)
	bot.Handle("/id", sendID)

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

func sendID(c tb.Context) error {
	if c.Chat() != nil {
		_ = c.Reply(fmt.Sprintf("Current Group ID: `%d`", c.Chat().ID), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
	} else {
		_ = c.Reply(fmt.Sprintf("Your Telegram ID: `%d`", c.Sender().ID), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
	}
	return nil
}
