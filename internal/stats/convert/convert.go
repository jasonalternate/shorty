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
func RepoSnapshotToModelSnapshot(in repo.Snapshot) model.Snapshot {
	return model.Snapshot{
		Slug: in.Slug,
		Count24Hours: in.Count24Hours,
		CountOneWeek: in.CountOneWeek,
		CountAllTime: in.CountAllTime,
	}
}
