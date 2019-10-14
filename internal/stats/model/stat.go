package model

import (
	"time"
)

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

