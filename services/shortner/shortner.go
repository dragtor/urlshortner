package shortner

import (
	"errors"
	"hash/crc32"
	"log"
	neturl "net/url"

	"github.com/catinello/base62"
	"github.com/dragtor/urlshortner/domain/url"
	"github.com/dragtor/urlshortner/domain/url/memory"
	"github.com/dragtor/urlshortner/domain/urlmetrics"
	metricsMemory "github.com/dragtor/urlshortner/domain/urlmetrics/memory"
)

type ShortnerConfiguration func(ss *ShortnerService) error

type ShortnerService struct {
	urlMeta    url.Repository
	urlMetrics urlmetrics.Repository
}

const (
	EVENT_CREATE_NEW_SHORTURL = "EVENT_CREATE_NEW_SHORTURL"
)

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
		ss.urlMetrics = metricsMemory.New()
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
	err = ss.UpdateMetrics(EVENT_CREATE_NEW_SHORTURL, sourceUrl)
	if err != nil {
		log.Println("Failed to update ")
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

func getDomainName(URL string) (string, error) {
	parsedURL, err := neturl.Parse(URL)
	if err != nil {
		return "", err
	}

	return parsedURL.Host, nil
}

func (ss *ShortnerService) GetMetrics(headCount int) ([]*urlmetrics.Metrics, error) {
	mts, err := ss.urlMetrics.GetTopCount(headCount)
	if err != nil {
		return nil, err
	}
	return mts, nil
}

func (ss *ShortnerService) UpdateMetrics(event, url string) error {
	host, err := getDomainName(url)
	if err != nil {
		log.Println(err)
		return err
	}
	var metrics *urlmetrics.Metrics
	_, err = ss.urlMetrics.GetMetrics(host)
	if err == metricsMemory.ErrDomainNotFound {
		metrics = urlmetrics.NewMetrics(host)
		metrics.IncrementCount()
		err = ss.urlMetrics.SetMetrics(host, metrics)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	metrics, err = ss.urlMetrics.GetMetrics(host)
	if err != nil {
		log.Println(err)
		return err
	}
	metrics.IncrementCount()
	err = ss.urlMetrics.SetMetrics(host, metrics)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
