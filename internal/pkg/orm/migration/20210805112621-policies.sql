-- +migrate Up
CREATE TABLE "policies" (
  "id" bigserial,
  "ptype" varchar(500) NOT NULL,
  "v0" varchar(500) NOT NULL,
  "v1" varchar(500) NOT NULL,
  "v2" varchar(500) NOT NULL,
  "v3" varchar(500) DEFAULT NULL,
  "v4" varchar(500) DEFAULT NULL,
  "v5" varchar(500) DEFAULT NULL
);

ALTER TABLE "policies"
  ADD CONSTRAINT policies_pkey PRIMARY KEY ("id");

CREATE INDEX policies_v2_idx ON "policies" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
WHERE
  "v3" IS NULL AND "v4" IS NULL AND "v5" IS NULL;

CREATE INDEX policies_v3_idx ON "policies" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
WHERE
  "v4" IS NULL AND "v5" IS NULL;

CREATE INDEX policies_v4_idx ON "policies" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
WHERE
  "v5" IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS policies_v4_idx;

DROP INDEX IF EXISTS policies_v3_idx;

DROP INDEX IF EXISTS policies_v2_idx;

ALTER TABLE "policies"
  DROP CONSTRAINT IF EXISTS policies_pkey;

DROP TABLE "policies";

