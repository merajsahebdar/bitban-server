-- +migrate Up
CREATE TABLE "casbin_rules" (
  "id" bigserial,
  "ptype" varchar(500) NOT NULL,
  "v0" varchar(500) NOT NULL,
  "v1" varchar(500) NOT NULL,
  "v2" varchar(500) NOT NULL,
  "v3" varchar(500) DEFAULT NULL,
  "v4" varchar(500) DEFAULT NULL,
  "v5" varchar(500) DEFAULT NULL
);

ALTER TABLE "casbin_rules"
  ADD CONSTRAINT casbin_rules_pkey PRIMARY KEY ("id");

CREATE INDEX casbin_rules_v2_idx ON "casbin_rules" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
WHERE
  "v3" IS NULL AND "v4" IS NULL AND "v5" IS NULL;

CREATE INDEX casbin_rules_v3_idx ON "casbin_rules" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
WHERE
  "v4" IS NULL AND "v5" IS NULL;

CREATE INDEX casbin_rules_v4_idx ON "casbin_rules" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
WHERE
  "v5" IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS casbin_rules_v4_idx;

DROP INDEX IF EXISTS casbin_rules_v3_idx;

DROP INDEX IF EXISTS casbin_rules_v2_idx;

ALTER TABLE "casbin_rules"
  DROP CONSTRAINT IF EXISTS casbin_rules_pkey;

DROP TABLE "casbin_rules";

