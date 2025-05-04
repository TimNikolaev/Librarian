package main

import (
	"context"
	"flag"
	tgClient "librarian/clients/telegram"
	"librarian/consumer/event_consumer"
	"librarian/events/telegram"
	"librarian/repository/sqlite"
	"log"
)

const (
	tgBotHost            = "api.telegram.org"
	filesRepositoryPath  = "files_repository"
	sqliteRepositoryPath = "data/sqlite/repository.db"
	batchSize            = 100
)

func main() {
	token := mustToken()

	// clients
	tgClient := tgClient.New(tgBotHost, token)

	//repository
	// filesRepository := files.New(filesRepositoryPath)
	sqliteRepository, err := sqlite.New(sqliteRepositoryPath)
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	if err = sqliteRepository.Init(context.TODO()); err != nil {
		log.Fatal("can't init repository:", err)
	}

	eventProcessor := telegram.New(tgClient, sqliteRepository)

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
