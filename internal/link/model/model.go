package model

type ShortLink struct {
	Slug string
	Destination string
}

func New(destination string) (ShortLink, error) {
	return ShortLink{}, nil
}
