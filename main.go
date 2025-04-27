package main

import (
	"flag"
	"librarian/clients/telegram"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	token := mustToken()

	_ = telegram.New(tgBotHost, token)

	//fetcher

	//processor

	//consumer
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
