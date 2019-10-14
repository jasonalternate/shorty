package service

import (
	"github.com/jasondeutsch/shorty/internal/stats/convert"
	"github.com/jasondeutsch/shorty/internal/stats/model"
	"github.com/jasondeutsch/shorty/internal/stats/repo"
	"time"
)

type Service interface {
	SaveView(model.View) error
	NewView(string, time.Time) model.View
	GetSnapshot(string) (model.Snapshot, error)
}

var _ Service = (*LocalService)(nil)

type LocalService struct {
	repo repo.Repository
}

func NewLocalService(repo repo.Repository) *LocalService {
	return &LocalService{repo: repo}
}

func (s *LocalService) SaveView(view model.View) error {
	return s.repo.SaveView(convert.ModelViewToRepoView(view))
}

func (s *LocalService) NewView(slug string, timeStamp time.Time) model.View{
	return model.View{
		Slug: slug,
		TimeStamp: timeStamp,
	}
}

func (s *LocalService) GetSnapshot(slug string) (model.Snapshot, error) {
	snapshot, err := s.repo.GetSnapshot(slug)

	return convert.RepoSnapshotToModelSnapshot(snapshot), err
}
