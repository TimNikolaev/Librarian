package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"librarian/pkg/e"
	"librarian/storage"
	"math/rand"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

const (
	defaultPerm = 0774
)

var ErrNoSavedPages = errors.New("no saved page")

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) error {
	const msgErr = "can't save page"

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return e.Wrap(msgErr, err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap(msgErr, err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap(msgErr, err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap(msgErr, err)
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (*storage.Page, error) {
	const msgErr = "can't pick random page"

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, e.Wrap(msgErr, err)
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	randFile := files[n]

	return s.decodePage(filepath.Join(fPath, randFile.Name()))
}

func (s *Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", fPath)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) IsExists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", fPath)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
