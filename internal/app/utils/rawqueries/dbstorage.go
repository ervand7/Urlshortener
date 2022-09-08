package rawqueries

const Set string = `insert into url("user_id", "short", "origin") values ($1, $2, $3);`
const Get string = `select "origin" from url where "short" = $1;`
const GetUserURLs string = `select "short", "origin" from url where "user_id" = $1;`
