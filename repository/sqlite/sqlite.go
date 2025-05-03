package sqlite

import (
	"context"
	"database/sql"
	"librarian/pkg/e"
	"librarian/repository"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(basePath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", basePath)
	if err != nil {
		return nil, e.Wrap("can't open database:", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database:", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	if _, err := r.db.ExecContext(ctx, q); err != nil {
		return e.Wrap("can't create table:", err)
	}

	return nil
}

func (r *Repository) Save(ctx context.Context, p *repository.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?,?)`

	if _, err := r.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (r *Repository) PickRandom(ctx context.Context, userName string) (*repository.Page, error) {
	var url string

	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	err := r.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap("can't pick random page:", err)
	}

	return &repository.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

func (r *Repository) Remove(ctx context.Context, p *repository.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	if _, err := r.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("can't remove page:", err)
	}

	return nil
}

func (r *Repository) IsExists(ctx context.Context, p *repository.Page) (bool, error) {
	var count int

	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	if err := r.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}
