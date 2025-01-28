package main

type ReqSurveillanceStation struct {
	ChatId  int64  `json:"chat_id"`
	Caption string `json:"caption"`
	Photo   string `json:"photo"`
	Silent  bool   `json:"silent"`
}
