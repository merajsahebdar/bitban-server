-- +migrate Up
CREATE TABLE "repositories" (
  "id" bigserial,
  "created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp with time zone NOT NULL DEFAULT NOW(),
  "removed_at" timestamp with time zone DEFAULT NULL,
  "address" varchar(250) NOT NULL,
  "domain_id" bigint DEFAULT NULL
);

ALTER TABLE "repositories"
  ADD CONSTRAINT repositories_pkey PRIMARY KEY ("id");

ALTER TABLE "repositories"
  ADD CONSTRAINT repositories_domain_fk FOREIGN KEY ("domain_id") REFERENCES "domains" ("id") ON DELETE CASCADE;

CREATE UNIQUE INDEX repositories_address_unq ON "repositories" ("domain_id", "address");

-- +migrate Down
DROP INDEX repositories_address_unq;

ALTER TABLE "repositories"
  DROP CONSTRAINT repositories_domain_fk;

ALTER TABLE "repositories"
  DROP CONSTRAINT repositories_pkey;

DROP TABLE "repositories";

