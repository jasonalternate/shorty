package convert

import (
	"github.com/jasondeutsch/shorty/internal/link/model"
	"github.com/jasondeutsch/shorty/internal/link/repo"
)

func ModelToRepo(in model.ShortLink) repo.ShortLink {
	return repo.ShortLink{
		Slug: in.Slug,
		Destination: in.Destination,
	}
}

func RepoToModel(in *repo.ShortLink) *model.ShortLink {
	if in == nil {
		return nil
	}
	return &model.ShortLink{
		Slug: in.Slug,
		Destination: in.Destination,
	}
}
