package url

import "errors"

var (
	ErrShortURLNotFound       = errors.New("short url not found")
	ErrSourceURLNotFound      = errors.New("source url not found")
	ErrShortURLAlreadyPresent = errors.New("url is already present")
)

type Repository interface {
	GetByShortUrl(shortUrl string) *UrlMeta
	Add(UrlMeta) error
	GetBySourceUrl(urlHash string) error
}
