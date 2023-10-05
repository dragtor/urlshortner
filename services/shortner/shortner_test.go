package shortner

import (
	"fmt"
	"testing"
)

func Test_Shortner(t *testing.T) {
	ss := &ShortnerService{}
	ss, err := NewShortnerService(WithMemoryUrlRepository(ss))
	if err != nil {
		t.Fatal(err)
	}
	sourceUrl := "https://www.notion.so/URL-Shortner-InfraCloud-a6aed1f0531349dbb4970c1ac62d4bb8"
	urlmeta, err := ss.CreateShortUrl(sourceUrl)
	if err != nil {
		t.Fatal(err)
	}
	shortUrl := urlmeta.GetShortUrl()
	urlmetaResponse, err := ss.GetSourceUrlForShortUrl(shortUrl)
	if err != nil {
		t.Fatal(err)
	}

	if urlmetaResponse.GetShortUrl() != urlmeta.GetShortUrl() {
		fmt.Println("db url : ", urlmetaResponse.GetShortUrl())
		fmt.Println("db url : ", urlmeta.GetShortUrl())
		t.Fatal("Url not matching ")
	}

	// try to add same source url
	urlmeta2, err := ss.CreateShortUrl(sourceUrl)
	if err != nil {
		t.Fatal(err)
	}
	url2shortUrl := urlmeta2.GetShortUrl()
	urlmetaResponse, err = ss.GetSourceUrlForShortUrl(url2shortUrl)
	if err != nil {
		t.Fatal(err)
	}

	if urlmetaResponse.GetShortUrl() != urlmeta.GetShortUrl() {
		fmt.Println("db url : ", urlmetaResponse.GetShortUrl())
		fmt.Println("db url : ", urlmeta.GetShortUrl())
		t.Fatal("Url not matching ")
	}

	sourceUrl = "https://www.livehindustan.com/national/story-chandrayaan-3-updates-vikram-and-pragyan-rover-may-awake-after-use-of-rtg-8790540.html"

	urlmeta2, err = ss.CreateShortUrl(sourceUrl)
	if err != nil {
		t.Fatal(err)
	}
	url2shortUrl = urlmeta2.GetShortUrl()
	urlmetaResponse, err = ss.GetSourceUrlForShortUrl(url2shortUrl)
	if err != nil {
		t.Fatal(err)
	}

	if urlmetaResponse.GetShortUrl() != urlmeta2.GetShortUrl() {
		fmt.Println("db url : ", urlmetaResponse.GetShortUrl())
		fmt.Println("db url : ", urlmeta.GetShortUrl())
		t.Fatal("Url not matching ")
	}
}

func Test_Metrics(t *testing.T) {

}
