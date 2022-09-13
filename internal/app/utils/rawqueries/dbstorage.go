package rawqueries

const (
	Set string = `insert into url("user_id", "short", "origin") values ($1, $2, $3);`

	CheckExists string = `select exists(select 1) from url where "origin" = $1;`

	GetOriginByShort string = `select "origin" from url where "short" = $1;`

	GetShortByOrigin string = `select "short" from url where "origin" = $1;`

	GetUserURLs string = `select "short", "origin" from url where "user_id" = $1;`
)
