-- +migrate Up
CREATE TABLE "users" (
  "domain_id" bigint NOT NULL,
  "domain_type" varchar(100) NOT NULL,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "password" varchar(250) DEFAULT NULL,
  "is_active" boolean NOT NULL,
  "is_banned" boolean NOT NULL
);

ALTER TABLE "users"
  ADD CONSTRAINT users_pkey PRIMARY KEY ("domain_id");

ALTER TABLE "users"
  ADD CONSTRAINT user_domain_check CHECK ("domain_type" = 'user');

ALTER TABLE "users"
  ADD CONSTRAINT user_domain_fk FOREIGN KEY ("domain_id", "domain_type") REFERENCES "domains" ("id", "type") ON DELETE CASCADE;

CREATE INDEX users_can_login_idx ON "users" ("domain_id", "domain_type", "password", "is_active", "is_banned")
WHERE
  "is_active" IS TRUE AND "is_banned" IS FALSE AND "removed_at" IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS users_can_login_idx;

ALTER TABLE "users"
  DROP CONSTRAINT user_domain_fk;

ALTER TABLE "users"
  DROP CONSTRAINT user_domain_check;

ALTER TABLE "users"
  DROP CONSTRAINT users_pkey;

DROP TABLE "users";

