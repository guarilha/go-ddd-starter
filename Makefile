OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
PKG ?= github.com/guarilha/go-ddd-starter/app
TERM=xterm-256color
CLICOLOR_FORCE=true
RICHGO_FORCE_COLOR=1
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_BUILD_TIME=$(shell date '+%Y-%m-%d__%I:%M:%S%p')
GO_BIN_PATH=$(shell go env GOPATH)/bin

define goBuild
	@echo "==> Go Building $2"
	@env GOOS=${OS} GOARCH=${ARCH} go build -v -o  build/$1 \
	-ldflags "-X main.BuildGitCommit=$(GIT_COMMIT) -X main.BuildTime=$(GIT_BUILD_TIME)" \
	${PKG}/$2
endef

.PHONY: compile
compile:
	$(call goBuild,service,"service")
	$(call goBuild,admin,"admin")

# ###########
# Setup
# ###########

.PHONY: install-moq
install-moq:
	@echo "==> Installing moq"
	@go install github.com/matryer/moq@latest

.PHONY: install-migration
install-migration:
	@echo "==> Installing migration"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: install-linters
install-linters:
	@echo "==> Installing linters"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install-test-fmt
install-test-fmt:
	@echo "==> Installing test formatter"
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest

.PHONY: install-gosec 
install-gosec:
	@echo "==> Installing gosec"
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

.PHONY: install-air
install-air:
	@echo "==> Installing Air (hot reload)"
	@go install github.com/air-verse/air@latest

.PHONY: install-sqlc 
install-sqlc:
	@echo "==> Installing sqlc"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


.PHONY: setup
setup: install-migration install-moq install-linters install-test-fmt install-gosec install-sqlc install-air
	@go mod tidy


# ###########
# Generate
# ###########

.PHONY: generate
generate: sqlc-generate 
	@echo "==> Running go generate"
	@go generate ./...

.PHONY: sqlc-generate
sqlc-generate:
	@echo "==> Generating sqlc code"
	@rm -f gateways/repository/*.go
	@sqlc generate


# ###########
# Lint
# ###########

.PHONY: lint 
lint:
	${GO_BIN_PATH}/golangci-lint run

# ###########
# GoSec 
# ###########

.PHONY: gosec 
gosec:
	${GO_BIN_PATH}/gosec -exclude-dir=gateways/repository ./...

# ###########
# Testing
# ###########

.PHONY: test-full
test-full:
	@go test -json -v -cover ./... 2>&1 | ${GO_BIN_PATH}/gotestfmt

.PHONY: test
test:
	@go test -json -v -short -cover ./... 2>&1 | ${GO_BIN_PATH}/gotestfmt

.PHONY: coverage
coverage:
	@go test -coverprofile=coverage.out ./... 2>&1 | ${GO_BIN_PATH}/gotestfmt
	@go tool cover -html=coverage.out

# ###########
# Migrations
# ###########

# Creates new migration up/down files in the 'migration' folder with the provided name.
.PHONY: migration/create
migration/create:
	@read -p "Enter migration name: " migration; \
	${GO_BIN_PATH}/migrate create -ext sql -dir ./gateways/repository/migrations/ "$$migration"

# Drop migration.
.PHONY: migration/drop
migration/drop:
	dsn="postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):5432/$(DATABASE_NAME)?sslmode=disable&search_path=public"; \
	${GO_BIN_PATH}/migrate -source file://gateways/repository/migrations -database $$dsn droprepository/migrations -seq $$migration

# Execute the migrations up to the most recent one. Needs the following environment variables:
# DATABASE_HOST: database url
# DATABASE_USER: database user
# DATABASE_PASSWORD: database password
# DATABASE_NAME: database name
.PHONY: migration/up
migration/up:
	dsn="postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):5432/$(DATABASE_NAME)?sslmode=disable&search_path=public"; \
	${GO_BIN_PATH}/migrate -source file://gateways/repository/migrations -database $$dsn up

# Rollback the migrations up to the oldest one. Needs the following environment variables:
# DATABASE_HOST: database url
# DATABASE_USER: database user
# DATABASE_PASSWORD: database password
# DATABASE_NAME: database name
.PHONY: migration/down
migration/down:
	dsn="postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):5432/$(DATABASE_NAME)?sslmode=disable&search_path=public"; \
	${GO_BIN_PATH}/migrate -source file://gateways/repository/migrations -database $$dsn drop
