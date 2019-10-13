package repo

import (
	"time"
)

type Repository interface {
	SaveView(View) error
}

type View struct {
	Slug string
	TimeStamp time.Time
}
