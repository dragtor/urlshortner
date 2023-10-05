package memory

import (
	"fmt"
	"testing"

	"github.com/dragtor/urlshortner/domain/urlmetrics"
)

func TestMemory_GetUrlMeta(t *testing.T) {
	mr := New()
	domain := "www.google.com"
	metrics := urlmetrics.NewMetrics(domain)
	err := mr.SetMetrics(domain, metrics)
	if err != nil {
		t.Fatal(err)
	}
	metrics.IncrementCount()
	metrics.IncrementCount()
	metrics.IncrementCount()
	err = mr.SetMetrics(domain, metrics)
	if err != nil {
		t.Fatal(err)
	}

	updateMetrics, err := mr.GetMetrics(domain)
	if err != nil {
		t.Fatal(err)
	}
	if updateMetrics.GetCount() == metrics.GetCount()+1 {
		t.Fatal(fmt.Sprintf("expected count %d , got %d", metrics.GetCount()+1, updateMetrics.GetCount()))
	}

	// domain : www.facebook.com
	domain = "www.facebook.com"
	metrics = urlmetrics.NewMetrics(domain)
	err = mr.SetMetrics(domain, metrics)
	if err != nil {
		t.Fatal(err)
	}
	metrics.IncrementCount()
	err = mr.SetMetrics(domain, metrics)
	if err != nil {
		t.Fatal(err)
	}

	res, err := mr.GetTopCount(3)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		fmt.Println(v.GetDomain(), v.GetCount())
	}
}
