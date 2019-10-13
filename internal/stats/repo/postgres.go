package repo

import "database/sql"

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) SaveView(view View)  error {
	q := `insert into stats_raw(slug, time_stamp) values($1, $2)`
	_, err := r.db.Exec(q, view.Slug, view.TimeStamp)
	return err
}