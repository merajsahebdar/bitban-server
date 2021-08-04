-- +migrate Up
CREATE TABLE "user_profiles" (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "name" varchar(250) NOT NULL,
  "meta" jsonb NOT NULL,
  "user_id" bigint DEFAULT NULL
);

ALTER TABLE "user_profiles"
  ADD CONSTRAINT user_profiles_pkey PRIMARY KEY ("id");

ALTER TABLE "user_profiles"
  ADD CONSTRAINT user_profiles_user_fk FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE INDEX user_profiles_user_idx ON "user_profiles" ("id", "name", "user_id");

-- +migrate Down
DROP INDEX IF EXISTS user_profiles_user_idx;

ALTER TABLE "user_profiles"
  DROP CONSTRAINT user_profiles_user_fk;

ALTER TABLE "user_profiles"
  DROP CONSTRAINT user_profiles_pkey;

DROP TABLE "user_profiles";

