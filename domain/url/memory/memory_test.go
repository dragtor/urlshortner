package memory

import (
	"fmt"
	"testing"

	"github.com/dragtor/urlshortner/domain/url"
)

func TestMemory_GetUrlMeta(t *testing.T) {
	type testCase struct {
		sourceUrl   string
		shortUrl    string
		expectedErr error
	}
	sourceUrl := "https://www.notion.so/URL-Shortner-InfraCloud-a6aed1f0531349dbb4970c1ac62d4bb8"
	shortUrl := "ntwerc"
	urlmeta, err := url.NewUrlMeta(sourceUrl, shortUrl)
	if err != nil {
		t.Fatal(err)
	}
	mr := New()
	err = mr.Add(urlmeta)
	if err != nil {
		t.Fatal(err)
	}

	urlmeta, err = mr.GetByShortUrl("ntwerc")
	if err != nil {
		t.Fatal(err)
	}
	if urlmeta.GetShortUrl() != shortUrl {
		t.Fatal("short url for given url not correct")
	}

	fmt.Println(urlmeta)
}
