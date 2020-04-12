package links

type Link struct {
	tableName struct{} `pg:"link"`
	ID        string   `json:"id"`
	LongURL   string   `json:"long_url"`
	ShortURL  string   `json:"short_url"`
}

type LinkRepository interface {
	Save(Link) (Link, error)
	FindLinkByID(string) (Link, error)
	FindLinkByShortURL(string) (Link, error)
}

type Request interface {
	Validate() bool
}

type Response interface {
	StatusCode() int
	Headers() map[string]string
	Body() interface{}
	Empty() bool
}
