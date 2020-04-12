package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mfinley3/links-v1/internal/links"
	linkRepoMock "github.com/mfinley3/links-v1/internal/links/postgres/mock"
)

func newService() Service {
	return New(linkRepoMock.New())
}

func TestSave(t *testing.T) {
	newLink := linkRepoMock.TestLinkOne
	newLink.ID = ""
	newLink.ShortURL = ""

	cases := []struct {
		desc        string
		link        links.Link
		expectedErr error
	}{
		{"Shoten link", newLink, nil},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			linkService := newService()
			link, err := linkService.Shorten(context.Background(), tc.link)
			require.Equal(t, tc.expectedErr, err)
			assert.Equal(t, linkRepoMock.TestLinkOne.LongURL, link.LongURL)
			assert.NotNil(t, link.ShortURL)
		})
	}
}

func TestFind(t *testing.T) {
	cases := []struct {
		desc         string
		linkID       string
		expectedLink links.Link
		expectedErr  error
	}{
		{"Find", linkRepoMock.TestLinkOne.ID, linkRepoMock.TestLinkOne, nil},
		{"Try to find a non-existent link", "fake link id", links.Link{}, ErrResourceNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			linkService := newService()
			link, err := linkService.Find(context.Background(), tc.linkID)
			require.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedLink, link)
		})
	}
}

func TestRedirect(t *testing.T) {
	cases := []struct {
		desc         string
		shortURL     string
		expectedLink links.Link
		expectedErr  error
	}{
		{"Redirect", linkRepoMock.TestLinkOne.ShortURL, linkRepoMock.TestLinkOne, nil},
		{"Try to find a non-existent link", "fake link id", links.Link{}, ErrResourceNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			linkService := newService()
			link, err := linkService.Redirect(context.Background(), tc.shortURL)
			require.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedLink, link)
		})
	}
}

func TestGenerateLinkNumber(t *testing.T) {
	cases := []struct {
		desc string
	}{
		{"Assert length is 8"},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			shortURL := generateShortLink()
			assert.Equal(t, 8, len(shortURL))
		})
	}
}
