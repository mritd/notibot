package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

var addr string
var token string
var botApi string
var botToken string
var recipient []int64

var rootCmd = &cobra.Command{
	Use:   "notibot",
	Short: "Telegram Notification Bot",
	Run: func(cmd *cobra.Command, args []string) {
		signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		bot, err := NewTelegram(botApi, botToken)
		if err != nil {
			logrus.Fatalf("failed to create telegram bot: %v", err)
		}

		app := fiber.New(fiber.Config{
			Prefork:       false,
			CaseSensitive: true,
			Immutable:     false,
			AppName:       "NotiBot",
		})

		go func() {
			app.Post("/message", func(ctx *fiber.Ctx) error {
				for _, r := range recipient {
					var err error
					if ctx.FormValue("markdown") != "true" {
						err = bot.SendMessage(ctx.FormValue("message"), r, true)
					} else {
						err = bot.SendMessage(ctx.FormValue("message"), r, false)
					}
					if err != nil {
						logrus.Errorf("failed to send: [%d] %v", r, err)
					}
				}
				return nil
			})
			app.Post("/file", func(ctx *fiber.Ctx) error {
				fh, err := ctx.FormFile("file")
				if err != nil {
					logrus.Errorf("failed to get file from request: %v", err)
					return err
				}
				f, err := fh.Open()
				if err != nil {
					logrus.Errorf("failed to get file from request: %v", err)
					return err
				}
				defer func() { _ = f.Close() }()
				for _, r := range recipient {
					err = bot.SendFile(f, fh.Filename, utils.GetMIME(fh.Filename), "", r)
					if err != nil {
						logrus.Errorf("failed to send: [%d] %v", r, err)
					}
				}
				return nil
			})
			app.Post("/image", func(ctx *fiber.Ctx) error {
				fh, err := ctx.FormFile("image")
				if err != nil {
					logrus.Errorf("failed to get image from request: %v", err)
					return err
				}
				f, err := fh.Open()
				if err != nil {
					logrus.Errorf("failed to get image from request: %v", err)
					return err
				}
				defer func() { _ = f.Close() }()
				for _, r := range recipient {
					err = bot.SendImage(f, "", r)
					if err != nil {
						logrus.Errorf("failed to send: [%d] %v", r, err)
					}
				}
				return nil
			})
			if err := app.Listen(addr); err != nil {
				logrus.Fatal(err)
			}
		}()

		logrus.Info("NotiBot Starting...")
		logrus.Infof("NotiBot API Token: %s", token)
		logrus.Infof("NotiBot Recipient: %v", recipient)

		<-signalCtx.Done()

		logrus.Info("NotiBot Shutdown!")
		_ = app.Shutdown()
	},
}

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	rootCmd.PersistentFlags().StringVarP(&addr, "listen", "l", "0.0.0.0:8080", "Server Listen Address")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", RandomString(32), "Server API Token")
	rootCmd.PersistentFlags().StringVarP(&botApi, "bot-api", "a", "https://api.telegram.org", "Telegram API Address")
	rootCmd.PersistentFlags().StringVarP(&botToken, "bot-token", "s", "", "Telegram Bot Token")
	rootCmd.PersistentFlags().Int64SliceVarP(&recipient, "recipient", "r", []int64{}, "Telegram Message Recipient")
	_ = rootCmd.MarkPersistentFlagRequired("bot-api")
	_ = rootCmd.MarkPersistentFlagRequired("bot-token")
	_ = rootCmd.MarkPersistentFlagRequired("recipient")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
