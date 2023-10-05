package url

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
)

type UrlMeta struct {
	url           string
	sourceURLHash string
	shortUrl      string
}

func generateSHA1(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	hashHex := fmt.Sprintf("%x", hashBytes)
	return hashHex
}

func NewUrlMeta(sourceUrl string, shortUrl string) (*UrlMeta, error) {
	if strings.TrimSpace(sourceUrl) == "" || strings.TrimSpace(shortUrl) == "" {
		return nil, errors.New("invalid values")
	}

	return &UrlMeta{
		url:           sourceUrl,
		sourceURLHash: generateSHA1(sourceUrl),
		shortUrl:      shortUrl,
	}, nil
}

func (u *UrlMeta) GetSourceURLHash() string {
	return u.sourceURLHash
}

func (u *UrlMeta) GetShortUrl() string {
	return u.shortUrl
}
