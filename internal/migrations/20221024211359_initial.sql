-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
create table if not exists "url"
(
    "id"      uuid                   DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" uuid          not null,
    "short"   varchar(30)   not null
        unique,
    "origin"  varchar(2000) not null
        unique,
    "active"  bool          not null default true
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table if exists "url";
