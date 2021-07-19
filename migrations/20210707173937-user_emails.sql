-- +migrate Up
CREATE TABLE user_emails (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "address" varchar(250) NOT NULL,
  "is_primary" boolean NOT NULL,
  "is_verified" boolean NOT NULL,
  "user_id" bigint DEFAULT NULL
);

ALTER TABLE "user_emails"
  ADD CONSTRAINT user_emails_pkey PRIMARY KEY ("id");

ALTER TABLE "user_emails"
  ADD CONSTRAINT user_emails_user_fk FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE UNIQUE INDEX user_emails_address_unq ON "user_emails" ("address");

CREATE INDEX user_emails_primary_idx ON "user_emails" ("id", "address", "user_id")
WHERE
  is_primary IS TRUE AND is_verified IS TRUE AND removed_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS user_emails_primary_idx;

DROP INDEX IF EXISTS user_emails_address_unq;

ALTER TABLE "user_emails"
  DROP CONSTRAINT user_emails_pkey;

DROP TABLE "user_emails";

