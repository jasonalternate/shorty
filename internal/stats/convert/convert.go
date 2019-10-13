package convert

import (
	"github.com/jasondeutsch/shorty/internal/stats/model"
	"github.com/jasondeutsch/shorty/internal/stats/repo"
)

func ModelViewToRepoView(in model.View) repo.View {
	return repo.View {
		Slug: in.Slug,
		TimeStamp: in.TimeStamp,
	}
}

func RepoViewToModelView(in repo.View) model.View {
	return model.View {
		Slug: in.Slug,
		TimeStamp: in.TimeStamp,
	}
}
