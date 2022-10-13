package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

var addr string
var authMode string
var accessToken string
var username string
var password string
var botApi string
var botToken string
var recipient []int64

var rootCmd = &cobra.Command{
	Use:   "notibot",
	Short: "Telegram Notification Bot",
	Run: func(cmd *cobra.Command, args []string) {
		signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		logrus.Info("NotiBot Starting...")
		logrus.Infof("NotiBot API User: %s", username)
		logrus.Infof("NotiBot API Password: %s", password)
		logrus.Infof("NotiBot Recipient: %v", recipient)

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
			switch authMode {
			case "none":
			case "user-password":
				app.Use(basicauth.New(basicauth.Config{
					Users: map[string]string{
						username: password,
					},
				}))
			case "access-token":
				app.Use(func(c *fiber.Ctx) error {
					if accessToken != c.Get("X-Auth") {
						c.Status(fiber.StatusForbidden)
						return c.Send([]byte("403 Forbidden"))
					}
					return c.Next()
				})
			default:
				logrus.Fatalf("unsupported auth mode: %v", authMode)
			}
			app.Post("/message", func(c *fiber.Ctx) error {
				for _, r := range recipient {
					var err error
					if c.FormValue("markdown") == "true" {
						err = bot.SendMessage(c.FormValue("message"), r, true)
					} else {
						err = bot.SendMessage(c.FormValue("message"), r, false)
					}
					if err != nil {
						logrus.Errorf("failed to send: [%d] %v", r, err)
						c.Status(fiber.StatusInternalServerError)
						_ = c.Send([]byte(err.Error()))
					}
				}
				return nil
			})
			app.Post("/file", func(c *fiber.Ctx) error {
				fh, err := c.FormFile("file")
				if err != nil {
					logrus.Errorf("failed to get file from request: %v", err)
					c.Status(fiber.StatusBadRequest)
					_ = c.Send([]byte(err.Error()))
					return err
				}
				f, err := fh.Open()
				if err != nil {
					logrus.Errorf("failed to get file from request: %v", err)
					c.Status(fiber.StatusInternalServerError)
					_ = c.Send([]byte(err.Error()))
					return err
				}
				defer func() { _ = f.Close() }()
				for _, r := range recipient {
					err = bot.SendFile(f, fh.Filename, utils.GetMIME(fh.Filename), "", r)
					if err != nil {
						logrus.Errorf("failed to send: [%d] %v", r, err)
						c.Status(fiber.StatusInternalServerError)
						_ = c.Send([]byte(err.Error()))
					}
				}
				return nil
			})
			app.Post("/image", func(c *fiber.Ctx) error {
				fh, err := c.FormFile("image")
				if err != nil {
					logrus.Errorf("failed to get image from request: %v", err)
					c.Status(fiber.StatusBadRequest)
					_ = c.Send([]byte(err.Error()))
					return err
				}
				f, err := fh.Open()
				if err != nil {
					logrus.Errorf("failed to get image from request: %v", err)
					c.Status(fiber.StatusInternalServerError)
					_ = c.Send([]byte(err.Error()))
					return err
				}
				defer func() { _ = f.Close() }()
				for _, r := range recipient {
					err = bot.SendImage(f, "", r)
					if err != nil {
						logrus.Errorf("failed to send: [%d] %v", r, err)
						c.Status(fiber.StatusInternalServerError)
						_ = c.Send([]byte(err.Error()))
					}
				}
				return nil
			})
			if err := app.Listen(addr); err != nil {
				logrus.Fatal(err)
			}
		}()

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
	rootCmd.PersistentFlags().StringVarP(&authMode, "auth-mode", "m", "access-token", "Server API Mode(access-token/user-password/none)")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "noti", "Server API Auth User")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", RandomString(16), "Server API Auth Password")
	rootCmd.PersistentFlags().StringVarP(&accessToken, "access-token", "t", RandomString(32), "Server API Auth AccessToken")
	rootCmd.PersistentFlags().StringVarP(&botApi, "bot-api", "a", "https://api.telegram.org", "Telegram API Address")
	rootCmd.PersistentFlags().StringVarP(&botToken, "bot-token", "s", "", "Telegram Bot Token")
	rootCmd.PersistentFlags().Int64SliceVarP(&recipient, "recipient", "r", []int64{}, "Telegram Message Recipient")
	_ = rootCmd.MarkPersistentFlagRequired("bot-token")
	_ = rootCmd.MarkPersistentFlagRequired("recipient")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
