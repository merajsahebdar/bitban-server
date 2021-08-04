-- +migrate Up
CREATE TABLE user_tokens (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "meta" jsonb NOT NULL,
  "user_id" bigint DEFAULT NULL
);

ALTER TABLE "user_tokens"
  ADD CONSTRAINT user_tokens_pkey PRIMARY KEY ("id");

ALTER TABLE "user_tokens"
  ADD CONSTRAINT user_tokens_user_fk FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE INDEX user_tokens_exist_idx ON "user_tokens" ("id", "user_id")
WHERE
  removed_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS user_tokens_primary_idx;

ALTER TABLE "user_tokens"
  DROP CONSTRAINT user_tokens_pkey;

DROP TABLE "user_tokens";

