package repo

import (
	"time"
)

type Repository interface {
	SaveView(View) error
	GetSnapshot(string) (Snapshot, error)
}

type View struct {
	Slug string
	TimeStamp time.Time
}

type Snapshot struct {
	Slug string
	Count24Hours int
	CountOneWeek int
	CountAllTime int
}
