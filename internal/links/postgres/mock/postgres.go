package mock

import (
	"github.com/mfinley3/links-v1/internal/links"
)

var _ links.LinkRepository = (*linkRepositoryMock)(nil)

var TestLinkOne = links.Link{
	ID:       "15175048-1e7c-4fe2-af53-e9b4182ddb04",
	LongURL:  "http://www.google.com",
	ShortURL: "Eg31buQ9",
}

type linkRepositoryMock struct {
	linksMap map[string]links.Link
}

func New() links.LinkRepository {
	var linksMap = make(map[string]links.Link)
	linksMap[TestLinkOne.ID] = TestLinkOne
	return &linkRepositoryMock{linksMap: linksMap}
}

func (lrm *linkRepositoryMock) Save(link links.Link) (links.Link, error) {
	lrm.linksMap[link.ID] = link
	return link, nil
}

func (lrm *linkRepositoryMock) FindLinkByID(id string) (links.Link, error) {
	link, ok := lrm.linksMap[id]
	if !ok {
		return links.Link{}, nil
	}
	return link, nil
}

func (lrm *linkRepositoryMock) FindLinkByShortURL(shortURL string) (links.Link, error) {
	for _, link := range lrm.linksMap {
		if link.ShortURL == shortURL {
			return link, nil
		}
	}
	return links.Link{}, nil
}
