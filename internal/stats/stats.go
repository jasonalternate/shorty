package stats

import (
	"database/sql"

	"github.com/jasondeutsch/shorty/internal/stats/repo"
	"github.com/jasondeutsch/shorty/internal/stats/service"
)

func NewStatsService(repo repo.Repository) service.Service {
	return service.NewLocalService(repo)
}

func NewPostgres(db *sql.DB) repo.Repository {
	return repo.NewPostgres(db)
}
