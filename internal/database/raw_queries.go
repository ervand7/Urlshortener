package database

const (
	Set string = `
		with cte as (
			insert into url ("user_id", "short", "origin")
				values ($1, $2, $3)
				on conflict ("origin") do nothing
				returning "short")
		select 'null'
		where exists(select 1 from cte)
		union all
		select "short"
		from url
		where "origin" = $3
		  and not exists(select 1 from cte);
`

	Get string = `select "origin", "active" from url where "short" = $1;`

	GetUserURLs string = `select "short", "origin" from url where "user_id" = $1;`

	DeleteURL string = `update url set "active" = false  where "short" = ANY($1)`
)
