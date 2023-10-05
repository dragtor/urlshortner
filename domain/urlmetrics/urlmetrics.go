package urlmetrics

type Metrics struct {
	domain string
	count  int
}

func NewMetrics(domain string) *Metrics {
	return &Metrics{
		domain: domain,
		count:  0,
	}
}

func (m *Metrics) IncrementCount() {
	m.count++
}

func (m *Metrics) GetCount() int {
	return m.count
}
func (m *Metrics) GetDomain() string {
	return m.domain
}
