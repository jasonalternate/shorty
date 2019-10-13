create table if not exists stats_raw (
    slug text not null,
    time_stamp timestamp not null default current_timestamp
);

create index if not exists statsRawSlugIdx on stats_raw(slug);

create table if not exists stats_snapshot (
    slug text not null primary key,
    count_24hours  integer not null default 0,
    count_one_week integer not null default 0,
    count_all_time integer not null default 0
);
