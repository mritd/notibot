package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func handleMessage(c *fiber.Ctx) error {
	markdown := c.FormValue("markdown") == "true"
	silent := c.FormValue("silent") == "true"

	for _, r := range recipient {
		if err := telegram.SendMessage(c.FormValue("message"), r, markdown, silent); err != nil {
			logrus.Errorf("failed to send: [%d] %v", r, err)
			c.Status(fiber.StatusInternalServerError)
			_ = c.Send([]byte(err.Error()))
		}
	}
	return nil
}

func handleFile(c *fiber.Ctx) error {
	silent := c.FormValue("silent") == "true"
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
		err = telegram.SendDocument(f, fh.Filename, "", r, silent)
		if err != nil {
			logrus.Errorf("failed to send: [%d] %v", r, err)
			c.Status(fiber.StatusInternalServerError)
			_ = c.Send([]byte(err.Error()))
		}
	}
	return nil
}

func handleImage(c *fiber.Ctx) error {
	silent := c.FormValue("silent") == "true"
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
		err = telegram.SendPhoto(f, "", r, silent)
		if err != nil {
			logrus.Errorf("failed to send: [%d] %v", r, err)
			c.Status(fiber.StatusInternalServerError)
			_ = c.Send([]byte(err.Error()))
		}
	}
	return nil
}

func handleSurveillanceStation(c *fiber.Ctx) error {
	var param ReqSurveillanceStation
	if err := c.BodyParser(&param); err != nil {
		return err
	}

	if param.Photo == "" {
		return fmt.Errorf("photo is empty")
	}

	bs, err := telegoutil.DownloadFile(param.Photo)
	if err != nil {
		return err
	}

	if param.ChatId != 0 {
		return telegram.SendPhoto(bytes.NewReader(bs), param.Caption, param.ChatId, param.Silent)
	} else {
		for _, r := range recipient {
			err := telegram.SendPhoto(bytes.NewReader(bs), param.Caption, param.ChatId, param.Silent)
			if err != nil {
				logrus.Errorf("failed to send: [%d] %v", r, err)
				c.Status(fiber.StatusInternalServerError)
				_ = c.Send([]byte(err.Error()))
			}
		}
		return nil
	}
}
