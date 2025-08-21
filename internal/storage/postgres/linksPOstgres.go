package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/k3mpton/shortner-project/pkg/conndb"
)

type Db struct {
	db *sql.DB
}

func NewStorage() *Db {
	DataB := conndb.Conn()
	return &Db{
		db: DataB,
	}
}

func (d *Db) LinkSave(
	ctx context.Context,
	originalLink string,
	shortLik string,
) error {
	const op = "postgres.LinkSave"

	query := `insert into Urls(original, short) values($1, $2)`
	row := d.db.QueryRow(query, originalLink, shortLik)
	if err := row.Err(); err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (d *Db) GetLink(
	ctx context.Context,
	shortLink string,
) (string, error) {
	const op = "postgres.GetLink"

	query := `select original from Urls
	 where short = $1
	`

	row := d.db.QueryRow(query, shortLink)
	if err := row.Err(); err != nil {
		return "", fmt.Errorf("%v: %v", op, err)
	}

	var origLink string
	if err := row.Scan(&origLink); err != nil {
		return "", fmt.Errorf("%v: %v", op, err)
	}

	return origLink, nil
}
