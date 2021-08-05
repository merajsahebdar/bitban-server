-- +migrate Up
CREATE TABLE "emails" (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "address" varchar(250) NOT NULL,
  "is_primary" boolean NOT NULL,
  "is_verified" boolean NOT NULL,
  "user_id" bigint DEFAULT NULL
);

ALTER TABLE "emails"
  ADD CONSTRAINT emails_pkey PRIMARY KEY ("id");

ALTER TABLE "emails"
  ADD CONSTRAINT emails_user_fk FOREIGN KEY ("user_id") REFERENCES "users" ("domain_id") ON DELETE CASCADE;

CREATE UNIQUE INDEX emails_address_unq ON "emails" ("address");

CREATE INDEX emails_primary_idx ON "emails" ("id", "address", "user_id")
WHERE
  is_primary IS TRUE AND is_verified IS TRUE AND removed_at IS NULL;

-- +migrate Down
DROP INDEX emails_primary_idx;

DROP INDEX emails_address_unq;

ALTER TABLE "emails"
  DROP CONSTRAINT emails_user_fk;

ALTER TABLE "emails"
  DROP CONSTRAINT emails_pkey;

DROP TABLE "emails";

