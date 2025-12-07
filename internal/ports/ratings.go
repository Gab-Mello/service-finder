package ports

type Ratings interface {
	AvgForProvider(providerID string) (avg float64, count int)
}
