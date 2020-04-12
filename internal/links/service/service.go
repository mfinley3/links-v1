package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/mfinley3/links-v1/internal/links"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
)

type Service interface {
	Shorten(context.Context, links.Link) (links.Link, error)
	Find(context.Context, string) (links.Link, error)
	Redirect(context.Context, string) (links.Link, error)
}

//How we can assert that a struct completely implements an interface
var _ Service = (*linkService)(nil)

type linkService struct {
	links links.LinkRepository
}

func New(lr links.LinkRepository) Service {
	return &linkService{
		links: lr,
	}
}

func (ls *linkService) Shorten(ctx context.Context, link links.Link) (links.Link, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return links.Link{}, err
	}
	link.ID = id.String()
	link.ShortURL = generateShortLink()

	link, err = ls.links.Save(link)
	if err != nil {
		return links.Link{}, err
	}

	return link, nil
}

func (ls *linkService) Find(ctx context.Context, id string) (links.Link, error) {
	link, err := ls.links.FindLinkByID(id)
	if err != nil {
		return links.Link{}, err
	}
	if link == (links.Link{}) {
		return links.Link{}, ErrResourceNotFound
	}
	return link, err
}

func (ls *linkService) Redirect(ctx context.Context, shortUrl string) (links.Link, error) {

	link, err := ls.links.FindLinkByShortURL(shortUrl)
	if err != nil {
		return links.Link{}, err
	}
	if link == (links.Link{}) {
		return links.Link{}, ErrResourceNotFound
	}
	return link, err
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShortLink() string {
	var validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	link := make([]byte, 8)
	for i := range link {
		link[i] = validChars[seededRand.Intn(len(validChars))]
	}
	return string(link)
}
