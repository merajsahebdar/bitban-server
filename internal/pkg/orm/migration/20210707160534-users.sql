-- +migrate Up
CREATE TABLE "users" (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "password" varchar(250) DEFAULT NULL,
  "is_active" boolean NOT NULL,
  "is_banned" boolean NOT NULL
);

ALTER TABLE "users"
  ADD CONSTRAINT users_pkey PRIMARY KEY ("id");

CREATE INDEX users_can_login_idx ON "users" ("id", "password", "is_active", "is_banned")
WHERE
  "is_active" IS TRUE AND "is_banned" IS FALSE AND "removed_at" IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS users_can_login_idx;

ALTER TABLE "users"
  DROP CONSTRAINT users_pkey;

DROP TABLE "users";

