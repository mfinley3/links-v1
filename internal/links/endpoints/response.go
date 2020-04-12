package endpoints

import (
	"net/http"

	"github.com/mfinley3/links-v1/internal/links"
)

var (
	_ links.Response = (*LinkResponse)(nil)
	_ links.Response = (*RedirectResponse)(nil)
)

// LINK RESPONSE

type LinkResponse struct {
	links.Link
	headers map[string]string
	created bool
}

func (r LinkResponse) StatusCode() int {
	if r.created {
		return http.StatusCreated
	}
	return http.StatusOK
}

// If the endpoints need them, they can return custom headers
func (r LinkResponse) Headers() map[string]string {
	return r.headers
}

func (r LinkResponse) Body() interface{} {
	return r
}

func (r LinkResponse) Empty() bool {
	return false
}

// REDIRECT RESPONSE

type RedirectResponse struct {
	headers map[string]string
}

func (r RedirectResponse) StatusCode() int {
	//Use StatusTemporaryRedirect because a web browser will remember (cache) a StatusPermanantRedirect and we'll miss out on data
	return http.StatusTemporaryRedirect
}

func (r RedirectResponse) Headers() map[string]string {
	return r.headers
}

func (r RedirectResponse) Body() interface{} {
	return r
}

func (r RedirectResponse) Empty() bool {
	return true
}
