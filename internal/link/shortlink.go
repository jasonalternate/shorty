package link

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jasondeutsch/shorty/internal/link/repo"
	"github.com/jasondeutsch/shorty/internal/link/service"
)

func NewShortLinkService(repo repo.Repository) service.Service {
	return service.NewLocalService(repo)
}

func NewMongo(db *mongo.Database) repo.Repository {
	return repo.NewMongo(db)
}

