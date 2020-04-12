package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics/influx"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/mfinley3/links-v1/internal/links"
	"github.com/mfinley3/links-v1/internal/links/endpoints"
	link "github.com/mfinley3/links-v1/internal/links/service"
	"github.com/mfinley3/links-v1/internal/links/transport"
	"github.com/mfinley3/links-v1/internal/middleware"
)

func Handler(ls link.Service, redirectCounter *influx.Counter) http.Handler {
	mux := chi.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
	}

	mux.Post("/links/v1/urls", kithttp.NewServer(
		endpoints.ShortenLink(ls),
		decodeShotenLinkRequest,
		encodeResponse,
		opts...,
	).ServeHTTP)

	mux.Get("/links/v1/urls/{link_id}", kithttp.NewServer(
		endpoints.GetLink(ls),
		decodeGetLinkRequest,
		encodeResponse,
		opts...,
	).ServeHTTP)

	mux.With(middleware.Timer(), middleware.RedirectCounter(redirectCounter)).Get("/{short_url}", kithttp.NewServer(
		endpoints.Redirect(ls),
		decodeRedirectRequest,
		encodeResponse,
		opts...,
	).ServeHTTP)

	return mux
}

func decodeShotenLinkRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var link links.Link
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		return nil, err
	}

	req := transport.ShortenReqest{
		Link: link,
	}

	return req, nil
}

func decodeGetLinkRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return transport.GetReqest{
		ID: chi.URLParam(r, "link_id"),
	}, nil
}

func decodeRedirectRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return transport.RedirectReqest{
		ShortURL: chi.URLParam(r, "short_url"),
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("content-type", "application/json")

	r, ok := resp.(links.Response)
	if !ok {
		return errors.New("error") //make proper error handling with encode error
	}

	for k, v := range r.Headers() {
		w.Header().Set(k, v)
	}
	w.WriteHeader(r.StatusCode())

	if r.Empty() {
		return nil
	}
	return json.NewEncoder(w).Encode(r.Body())
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	errorResponse := struct {
		Message string
	}{
		Message: err.Error(),
	}

	switch err {
	case endpoints.ErrInvalidURL, endpoints.ErrInvalidID, endpoints.ErrInvalidShortURL:
		w.WriteHeader(http.StatusBadRequest)
	case link.ErrResourceNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(errorResponse)

}
