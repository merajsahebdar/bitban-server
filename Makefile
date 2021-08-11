# The directory where this Makefile lives in.
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# Migrations Directory
MIGRATIONS_DIR:=$(ROOT_DIR)/internal/pkg/orm/migration

# Phony!
.PHONY: build

# Build Folio Binary
build: clean pkger orm
	go build -o bitban

# Clean-Up
clean:
	rm -f bitban

# Deps
deps: export GO111MODULE=off
deps:
	go get github.com/joho/godotenv/cmd/godotenv

# Static Packing
pkger:
	rm -f pkged.go
	pkger

# Generate GraphQL Server and Schema
schema:
	rm -f internal/schema/server.go
	gqlgen generate

# Generate Protobuf
proto:
	find . -type f -name '*.pb.go' -delete
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(shell find . -type f -name '*.proto')

# Create a new migration template.
migration:
	@while [ -z "$$MIGRATION_NAME" ]; do \
        read -r -p "Give the migration a name: " MIGRATION_NAME;\
    done \
	&& TIMESTAMP=`date -u +%Y%m%d%H%M%S` \
	&& echo "-- +migrate Up\n\n-- +migrate Down" > "${MIGRATIONS_DIR}/$$TIMESTAMP-$$MIGRATION_NAME.sql"
