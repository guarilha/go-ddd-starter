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
	@echo "==> Go mod tidy"
	@go mod tidy
	$(call goBuild,service,"service")
	$(call goBuild,admin,"admin")

# ###########
# Setup
# ###########

.PHONY: enable-mise-experimental
enable-mise-experimental:
	@echo "==> Enabling mise experimental features"
	@mise settings experimental=true

.PHONY: install-mise-tools
install-mise-tools: enable-mise-experimental
	@echo "==> Installing mise dev tools"
	@mise install

# Workaround because mise currently don't support tags for the go backend
.PHONY: install-migration
install-migration: install-mise-tools
	@echo "==> Installing migration"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: setup-linter
setup-linter:
	@echo "==> Setting up custom linters"
	@cd internal/linters/domain_func_signature && go mod tidy
	@cd internal/linters/layer_imports && go mod tidy
	@rm -f internal/linters/domain_func_signature/plugin/plugin.so
	@rm -f internal/linters/layer_imports/plugin/plugin.so
	@make build-linter-plugin

.PHONY: setup
setup: install-migration setup-linter
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

.PHONY: build-linter-plugin
build-linter-plugin:
	@echo "==> Building domain function signature linter plugin"
	@cd internal/linters/domain_func_signature && CGO_ENABLED=1 go build -buildmode=plugin -o plugin/plugin.so plugin/plugin.go
	@echo "==> Building layer imports linter plugin"
	@cd internal/linters/layer_imports && CGO_ENABLED=1 go build -buildmode=plugin -o plugin/plugin.so plugin/plugin.go

.PHONY: lint
lint: build-linter-plugin
	@echo "==> Running linter"
	@mise x -- golangci-lint cache clean
	@mise x -- golangci-lint run --config=.golangci.yml ./...

# ###########
# GoSec 
# ###########

.PHONY: gosec 
gosec:
	mise x -- gosec -exclude-dir=gateways/pg ./...

# ###########
# Testing
# ###########

.PHONY: test-full
test-full:
	@go test -json -v -cover ./... 2>&1 | gotestfmt

.PHONY: test
test:
	@go test -json -v -short -cover ./... 2>&1 | gotestfmt

.PHONY: coverage
coverage:
	@go test -coverprofile=coverage.out ./... 2>&1 | gotestfmt
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
