package telegram

import "librarian/clients/telegram"

type Service struct {
	tg     *telegram.Client
	offset int
	//storage
}

func New(client *telegram.Client) *Service {
	return &Service{}
}
