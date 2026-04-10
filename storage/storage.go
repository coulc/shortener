package storage

type Storage interface {
	Save(*URLMapping) error
	Get(string) (*URLMapping,error)
	IncrementVisit(string) error
	Delete(string) error 
}

type URLMapping struct {
	ShortCode string
	LongURL string
	CreatedAt int64
	VisitCount int
}
