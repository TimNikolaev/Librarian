package main

import (
	"flag"
	tgClient "librarian/clients/telegram"
	event_consumer "librarian/consumer/event-consumer"
	"librarian/events/telegram"
	"librarian/repository/files"
	"log"
)

const (
	tgBotHost      = "api.telegram.org"
	repositoryPath = "repository"
	batchSize      = 100
)

func main() {
	token := mustToken()

	tgClient := tgClient.New(tgBotHost, token)
	filesRepository := files.New(repositoryPath)

	eventProcessor := telegram.New(tgClient, filesRepository)

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped")
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
