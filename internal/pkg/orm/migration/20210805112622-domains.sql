-- +migrate Up
CREATE TABLE "domains" (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "type" varchar(100) NOT NULL,
  "name" varchar(250) NOT NULL,
  "address" varchar(250) NOT NULL,
  "meta" jsonb NOT NULL
);

ALTER TABLE "domains"
  ADD CONSTRAINT domains_pkey PRIMARY KEY ("id");

ALTER TABLE "domains"
  ADD CONSTRAINT domains_type_check CHECK ("type" IN ('user', 'organization'));

CREATE UNIQUE INDEX domains_type_unq ON "domains" ("id", "type");

CREATE UNIQUE INDEX domains_address_unq ON "domains" ("address");

-- +migrate Down
DROP INDEX domains_address_unq;

DROP INDEX domains_type_unq;

ALTER TABLE "domains"
  DROP CONSTRAINT domains_type_check;

ALTER TABLE "domains"
  DROP CONSTRAINT domains_pkey;

DROP TABLE "domains";

