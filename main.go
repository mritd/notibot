package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/keyauth/v2"
	_ "github.com/mritd/logrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var recipient []int64
var telegram *Telegram

var rootCmd = &cobra.Command{
	Use:   "notibot",
	Short: "Telegram Notification Bot",
	Run: func(cmd *cobra.Command, args []string) {
		signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		logrus.Info("NotiBot Starting...")
		logrus.Infof("NotiBot Auth Mode: %s", viper.GetString("auth-mode"))
		logrus.Infof("NotiBot API User: %s", viper.GetString("username"))
		logrus.Infof("NotiBot API Password: %s", viper.GetString("password"))
		logrus.Infof("NotiBot API Access Token: %s", viper.GetString("access-token"))
		logrus.Infof("NotiBot Telegram Recipient: %v", viper.GetString("recipient"))

		if viper.GetString("recipient") == "" {
			logrus.Fatal("telegram recipient cannot be empty")
		}
		ss := strings.Split(viper.GetString("recipient"), ",")
		for _, s := range ss {
			id, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				logrus.Fatalf("failed to parse recipient: %v", err)
			}
			recipient = append(recipient, id)
		}

		bot, err := NewTelegram(viper.GetString("bot-api"), viper.GetString("bot-token"))
		if err != nil {
			logrus.Fatalf("failed to create telegram bot: %v", err)
		}
		telegram = bot

		app := fiber.New(fiber.Config{
			Prefork:       false,
			CaseSensitive: true,
			Immutable:     false,
			AppName:       "NotiBot",
		})

		go func() {
			switch viper.GetString("auth-mode") {
			case "none":
			case "basicauth":
				app.Use(basicauth.New(basicauth.Config{
					Users: map[string]string{
						viper.GetString("username"): viper.GetString("password"),
					},
				}))
			case "keyauth":
				app.Use(keyauth.New(keyauth.Config{
					KeyLookup:  "header:Authorization",
					AuthScheme: "Bearer",
					Validator: func(ctx *fiber.Ctx, s string) (bool, error) {
						return s == viper.GetString("access-token"), nil
					},
				}))
			default:
				logrus.Fatalf("unsupported auth mode: %v", viper.GetString("auth-mode"))
			}

			app.Post("/message", handleMessage)
			app.Post("/file", handleFile)
			app.Post("/image", handleImage)
			app.Post("/surveillance", handleSurveillanceStation)

			if err := app.Listen(viper.GetString("listen")); err != nil {
				logrus.Fatal(err)
			}
		}()

		<-signalCtx.Done()

		logrus.Info("NotiBot Shutdown!")
		_ = app.Shutdown()
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("listen", "l", "", "Server Listen Address")
	rootCmd.PersistentFlags().StringP("auth-mode", "m", "", "Server API Auth Mode(basicauth/keyauth/none)")
	rootCmd.PersistentFlags().StringP("username", "u", "", "Server API Basic Auth User")
	rootCmd.PersistentFlags().StringP("password", "p", "", "Server API Basic Auth Password")
	rootCmd.PersistentFlags().StringP("access-token", "t", "", "Server API AccessToken")
	rootCmd.PersistentFlags().StringP("bot-api", "a", "https://api.telegram.org", "Telegram API Address")
	rootCmd.PersistentFlags().StringP("bot-token", "s", "", "Telegram Bot Token")
	rootCmd.PersistentFlags().StringP("recipient", "r", "", "Telegram Message Recipient")

	_ = viper.BindEnv("listen", "NOTI_LISTEN_ADDR", "NOTI_LISTEN")
	_ = viper.BindEnv("auth-mode", "NOTI_AUTH_MODE")
	_ = viper.BindEnv("username", "NOTI_AUTH_USERNAME", "NOTI_USERNAME")
	_ = viper.BindEnv("password", "NOTI_AUTH_PASSWORD", "NOTI_PASSWORD")
	_ = viper.BindEnv("access-token", "NOTI_AUTH_TOKEN", "NOTI_TOKEN")
	_ = viper.BindEnv("bot-api", "NOTI_TELEGRAM_API")
	_ = viper.BindEnv("bot-token", "NOTI_TELEGRAM_TOKEN")
	_ = viper.BindEnv("recipient", "NOTI_TELEGRAM_RECIPIENT")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		logrus.Fatal(err)
	}

	if viper.GetString("listen") == "" {
		viper.Set("listen", "0.0.0.0:8080")
	}
	if viper.GetString("auth-mode") == "" {
		viper.Set("auth-mode", "keyauth")
	}
	if viper.GetString("access-token") == "" {
		viper.Set("access-token", RandomString(32))
	}
	if viper.GetString("username") == "" {
		viper.Set("username", "noti")
	}
	if viper.GetString("password") == "" {
		viper.Set("password", RandomString(16))
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
