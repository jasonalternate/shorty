package repo

import (
	"database/sql"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) SaveView(view View) error {
	q := `insert into stats_raw(slug, time_stamp) values($1, $2)`
	_, err := r.db.Exec(q, view.Slug, view.TimeStamp)
	return err
}

func (r *Postgres) GetSnapshot(slug string) (Snapshot, error) {
	q := `
select 
slug,
count(slug) filter (where time_stamp >= now() - interval '1 day') as count_24_hours, 
count(slug) filter (where time_stamp >= now() - interval '1 week') as count_one_week,
count(*) as count_all_time
 from stats_raw
 where slug = $1
 group by slug;
`
	var snapshot Snapshot
	err := r.db.QueryRow(q, slug).Scan(&snapshot.Slug, &snapshot.Count24Hours, &snapshot.CountOneWeek, &snapshot.CountAllTime)

	return snapshot, err
}
