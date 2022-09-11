package rawqueries

const CreateAll string = `
		create table if not exists "url"
		(
			"id"      uuid DEFAULT gen_random_uuid() PRIMARY KEY,
			"user_id" uuid          not null,
			"short"   varchar(30)   not null unique,
			"origin"  varchar(2000) not null
		);
		
		create index if not exists "ix_user_id"
			on "public"."url" using btree ("user_id");

		create index if not exists "ix_short"
					on "public"."url" using btree ("short");
`
