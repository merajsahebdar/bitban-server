-- +migrate Up
CREATE TABLE "tokens" (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "meta" jsonb NOT NULL,
  "user_id" bigint DEFAULT NULL
);

ALTER TABLE "tokens"
  ADD CONSTRAINT tokens_pkey PRIMARY KEY ("id");

ALTER TABLE "tokens"
  ADD CONSTRAINT tokens_user_fk FOREIGN KEY ("user_id") REFERENCES "users" ("domain_id") ON DELETE CASCADE;

CREATE INDEX tokens_exist_idx ON "tokens" ("id", "user_id")
WHERE
  removed_at IS NULL;

-- +migrate Down
DROP INDEX tokens_exist_idx;

ALTER TABLE "tokens"
  DROP CONSTRAINT tokens_user_fk;

ALTER TABLE "tokens"
  DROP CONSTRAINT tokens_pkey;

DROP TABLE "tokens";

