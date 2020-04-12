package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	link "github.com/mfinley3/links-v1/internal/links/service"
	"github.com/mfinley3/links-v1/internal/links/transport"
)

var (
	ErrInvalidURL      = errors.New("invalid url")
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidShortURL = errors.New("invalid short url")
)

func ShortenLink(ls link.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.ShortenReqest)
		if ok := req.Validate(); !ok {
			return nil, ErrInvalidURL
		}

		link, err := ls.Shorten(ctx, req.Link)
		if err != nil {
			return nil, err
		}

		return LinkResponse{
			Link:    link,
			created: true,
		}, nil
	}
}

func GetLink(ls link.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.GetReqest)
		if ok := req.Validate(); !ok {
			return nil, ErrInvalidID
		}

		link, err := ls.Find(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return LinkResponse{
			Link: link,
		}, nil
	}
}

func Redirect(ls link.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.RedirectReqest)
		if ok := req.Validate(); !ok {
			return nil, ErrInvalidShortURL
		}

		link, err := ls.Redirect(ctx, req.ShortURL)
		if err != nil {
			return nil, err
		}

		return RedirectResponse{
			headers: map[string]string{
				"location": link.LongURL,
			},
		}, nil
	}
}
