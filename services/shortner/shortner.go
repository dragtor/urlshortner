package shortner

import (
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/catinello/base62"
	"github.com/dragtor/urlshortner/domain/url"
	"github.com/dragtor/urlshortner/domain/url/memory"
)

type ShortnerConfiguration func(ss *ShortnerService) error

type ShortnerService struct {
	urlMeta url.Repository
}

func NewShortnerService(cfgs ...ShortnerConfiguration) (*ShortnerService, error) {
	ss := &ShortnerService{}
	for _, cfg := range cfgs {
		err := cfg(ss)
		if err != nil {
			return nil, err
		}
	}
	return ss, nil
}

func WithMemoryUrlRepository(ss *ShortnerService) ShortnerConfiguration {
	return func(ss *ShortnerService) error {
		ss.urlMeta = memory.New()
		return nil
	}
}

const (
	EmptyString = ""
)

func (ss *ShortnerService) validateSourceUrl(sourceURL string) error {
	if sourceURL == EmptyString {
		return errors.New("invalid input")
	}
	return nil
}

func generateCRC32Encoding(input string) uint32 {
	crc32Hash := crc32.NewIEEE()
	crc32Hash.Write([]byte(input))
	checksum := crc32Hash.Sum32()
	return checksum
}

func (ss *ShortnerService) CreateShortUrl(sourceUrl string) (*url.UrlMeta, error) {
	err := ss.validateSourceUrl(sourceUrl)
	if err != nil {
		return nil, err
	}
	crc32Encoding := generateCRC32Encoding(sourceUrl)
	shortUrl := base62.Encode(int(crc32Encoding))
	fmt.Println(shortUrl)
	urlmeta, err := url.NewUrlMeta(sourceUrl, shortUrl)
	if err != nil {
		return nil, err
	}
	err = ss.urlMeta.Add(urlmeta)
	if err == url.ErrShortURLAlreadyPresent {
		urlmeta, err = ss.urlMeta.GetBySourceUrl(urlmeta.GetSourceURLHash())
		if err != nil {
			return nil, err
		}
		return urlmeta, nil
	}
	if err != nil {
		return nil, err
	}
	return urlmeta, nil
}

func (ss *ShortnerService) validateShortUrl(url string) error {
	if url == EmptyString {
		return errors.New("invalid input")
	}
	return nil
}

func (ss *ShortnerService) GetSourceUrlForShortUrl(shortUrl string) (*url.UrlMeta, error) {
	err := ss.validateShortUrl(shortUrl)
	if err != nil {
		return nil, errors.New("invalid input")
	}
	urlmeta, err := ss.urlMeta.GetByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return urlmeta, nil
}
