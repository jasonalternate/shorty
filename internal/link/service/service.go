package service

import (
	"errors"
	"github.com/jasondeutsch/shorty/internal/keygen"
	"github.com/jasondeutsch/shorty/internal/link/convert"
	"github.com/jasondeutsch/shorty/internal/link/model"
	"github.com/jasondeutsch/shorty/internal/link/repo"
	"net/url"
)

type Service interface {
	Create(destination string) (*model.ShortLink, error)
	ReadOne(slug string) (*model.ShortLink, error)
}

var _ Service = (*LocalService)(nil)

type LocalService struct {
	repo repo.Repository
}

func NewLocalService(repo repo.Repository) *LocalService {
	return &LocalService{repo: repo}
}


func (s *LocalService) Create(destination string) (*model.ShortLink, error) {
	if _, err := url.ParseRequestURI(destination); err != nil {
		return nil, errors.New("invalid destination address")
	}


	slug, err := s.shortenUrlWithRetry(destination, 10)
	if err != nil {
		return nil, errors.New("unable to create shortlink")
	}

	link := model.ShortLink{
		Slug: slug,
		Destination: destination,
	}

	_, err = s.repo.Create(convert.ModelToRepo(link))

	return &link, err
}

func (s *LocalService) ReadOne(slug string) (*model.ShortLink, error) {
	link, err := s.repo.ReadOne(slug)

	return convert.RepoToModel(link), err
}

func (s *LocalService) shortenUrlWithRetry(destination string, retries int) (string, error){
	var slug string
	var err error
	existingLink := &model.ShortLink{}

	for i := 0; i < 10 && existingLink != nil; i++ {
		slug = keygen.KeyGenerator{}.Generate(7)

		existingLink, err = s.ReadOne(slug)
		if err != nil {
			return "", err
		}

	}
	return slug, nil
}
