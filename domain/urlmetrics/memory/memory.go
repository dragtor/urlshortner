package memory

import (
	"errors"
	"sort"
	"sync"

	"github.com/dragtor/urlshortner/domain/urlmetrics"
)

type MemoryRepository struct {
	domainMetrics map[string]*urlmetrics.Metrics
	sync.Mutex
}

var (
	ErrDomainNotFound = errors.New("domain not found")
	ErrEmptyValue     = errors.New("empty value")
)

func New() *MemoryRepository {
	return &MemoryRepository{
		domainMetrics: make(map[string]*urlmetrics.Metrics),
	}
}

func (mr *MemoryRepository) GetMetrics(domain string) (*urlmetrics.Metrics, error) {
	if metrics, ok := mr.domainMetrics[domain]; ok {
		return metrics, nil
	}
	return nil, ErrDomainNotFound
}

func (mr *MemoryRepository) SetMetrics(domain string, mtc *urlmetrics.Metrics) error {
	mr.Lock()
	defer mr.Unlock()
	if mtc == nil {
		return ErrEmptyValue
	}
	mr.domainMetrics[domain] = mtc
	return nil
}

type MetricsSlice []*urlmetrics.Metrics

func (ms MetricsSlice) Len() int {
	return len(ms)
}

func (ms MetricsSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms MetricsSlice) Less(i, j int) bool {
	return ms[i].GetCount() > ms[j].GetCount()
}

func (mr *MemoryRepository) GetTopCount(headCount int) ([]*urlmetrics.Metrics, error) {
	var metrics []*urlmetrics.Metrics
	for _, mt := range mr.domainMetrics {
		metrics = append(metrics, mt)
	}
	sort.Sort(MetricsSlice(metrics))
	if len(metrics) < headCount {
		return metrics, nil
	}
	return metrics[:headCount], nil
}
