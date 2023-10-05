package url

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"
)

type UrlMeta struct {
	url           string
	sourceURLHash string
	shortUrl      string
}

func GenerateSHA1(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	sha1Hex := hex.EncodeToString(hashBytes)
	return sha1Hex
}

func NewUrlMeta(sourceUrl string, shortUrl string) (*UrlMeta, error) {
	if strings.TrimSpace(sourceUrl) == "" || strings.TrimSpace(shortUrl) == "" {
		return nil, errors.New("invalid values")
	}

	return &UrlMeta{
		url:           sourceUrl,
		sourceURLHash: GenerateSHA1(sourceUrl),
		shortUrl:      shortUrl,
	}, nil
}

func (u *UrlMeta) GetSourceURLHash() string {
	return u.sourceURLHash
}

func (u *UrlMeta) GetShortUrl() string {
	return u.shortUrl
}

func (u *UrlMeta) GetSourceUrl() string {
	return u.url
}
