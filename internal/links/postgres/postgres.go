package postgres

import (
	"github.com/go-pg/pg/v9"

	"github.com/mfinley3/links-v1/internal/links"
)

var _ links.LinkRepository = (*linkRepository)(nil)

type linkRepository struct {
	db *pg.DB
}

func New(db *pg.DB) links.LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (lr *linkRepository) Save(link links.Link) (links.Link, error) {
	insertStatement := `INSERT INTO link (id, long_url, short_url)
						VALUES (? ,? ,?)
						RETURNING *`
	_, err := lr.db.Query(&link, insertStatement, link.ID, link.LongURL, link.ShortURL)
	return link, err
}

func (lr *linkRepository) FindLinkByID(id string) (links.Link, error) {
	var link links.Link
	selectStatement := `SELECT * FROM link WHERE id = ?;`
	_, err := lr.db.Query(&link, selectStatement, id)
	return link, err
}

func (lr *linkRepository) FindLinkByShortURL(url string) (links.Link, error) {
	var link links.Link
	selectStatement := `SELECT * FROM link WHERE short_url = ?;`
	_, err := lr.db.Query(&link, selectStatement, url)
	return link, err
}
