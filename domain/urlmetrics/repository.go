package urlmetrics

type Repository interface {
	GetMetrics(domain string) (*Metrics, error)
	SetMetrics(domain string, metrics *Metrics) error
	GetTopCount(headCount int) ([]*Metrics, error)
}
