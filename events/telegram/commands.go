package telegram

import (
	"errors"
	"librarian/pkg/e"
	"librarian/repository"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *EventProcessor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' form '%s'", text, userName)

	if isAddCmd(text) {
		return p.savePage(chatID, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *EventProcessor) savePage(chatID int, pageURL string, userName string) error {
	const msgErr = "can't do command: save page"

	page := &repository.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExists, err := p.repository.IsExists(page)
	if err != nil {
		return e.Wrap(msgErr, err)
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.repository.Save(page); err != nil {
		return e.Wrap(msgErr, err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap(msgErr, err)
	}

	return nil
}

func (p *EventProcessor) sendRandom(chatID int, userName string) error {
	const msgErr = "can't do command: can't send random"

	page, err := p.repository.PickRandom(userName)
	if err != nil && !errors.Is(err, repository.ErrNoSavedPages) {
		return e.Wrap(msgErr, err)
	}
	if errors.Is(err, repository.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap(msgErr, err)
	}

	return p.repository.Remove(page)
}

func (p *EventProcessor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *EventProcessor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
