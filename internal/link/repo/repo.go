package repo

type Repository interface {
	Create(ShortLink) (*ShortLink, error)
	ReadOne(string) (*ShortLink, error)
}

type ShortLink struct {
	Slug string `bson:"slug"`
	Destination string `bson:"destination"`
}
