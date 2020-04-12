package transport

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/mfinley3/links-v1/internal/links"
)

type ShortenReqest struct {
	Link links.Link
}

type GetReqest struct {
	id string
}

type RedirectReqest struct {
	ShortURL string
}

func (r ShortenReqest) Validate() bool {
	_, err := url.ParseRequestURI(r.Link.LongURL)
	if err != nil {
		return false
	}
	return true
}

func (r GetReqest) Validate() bool {
	_, err := uuid.Parse(r.id)
	if err != nil {
		return false
	}
	return true
}

func (r RedirectReqest) Validate() bool {
	return len(r.ShortURL) == 8
}
