package memory

import (
	"log"
	"sync"

	"github.com/dragtor/urlshortner/domain/url"
)

type MemoryRepository struct {
	urlMetaStore  map[string]*url.UrlMeta
	shortUrlStore map[string]*url.UrlMeta
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		urlMetaStore:  make(map[string]*url.UrlMeta),
		shortUrlStore: make(map[string]*url.UrlMeta),
	}
}

func (mr *MemoryRepository) GetByShortUrl(shortUrl string) (*url.UrlMeta, error) {
	if urlmeta, ok := mr.shortUrlStore[shortUrl]; ok {
		return urlmeta, nil
	}
	return nil, url.ErrShortURLNotFound
}
func (mr *MemoryRepository) Add(urlmeta *url.UrlMeta) error {
	log.Println("INFO: inmemory repository adding urlmeta for ", urlmeta.GetSourceUrl())
	if _, ok := mr.urlMetaStore[urlmeta.GetSourceURLHash()]; ok {
		return url.ErrShortURLAlreadyPresent
	}
	mr.Lock()
	defer mr.Unlock()
	log.Println("INFO: Compeleted inmemory repository adding urlmeta for ", urlmeta.GetSourceUrl())
	mr.urlMetaStore[urlmeta.GetSourceURLHash()] = urlmeta
	mr.shortUrlStore[urlmeta.GetShortUrl()] = urlmeta
	return nil
}
func (mr *MemoryRepository) GetBySourceUrl(urlHash string) (*url.UrlMeta, error) {
	log.Println("INFO: inmemory repository fetching source url by urlhash ", urlHash)
	if urlmeta, ok := mr.urlMetaStore[urlHash]; ok {
		return urlmeta, nil
	}
	return nil, url.ErrSourceURLNotFound
}
