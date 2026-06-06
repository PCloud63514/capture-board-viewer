PROJECT_NAME := capture-board-selector
PKG_PATH_BATCH := ./
SENTRY_DSN ?=
VERSION ?= v1.1.2

LDFLAGS := -X 'capture-board-selector/internal/captureboard.SentryDSN=$(SENTRY_DSN)' \
           -X 'capture-board-selector/internal/captureboard.Version=$(VERSION)'

verify:
	go mod tidy
	go mod vendor

ci:
	go mod tidy
	go mod vendor
	go build -ldflags "$(LDFLAGS)" -o $(PROJECT_NAME).exe $(PKG_PATH_BATCH)

run:
	$(MAKE) ci
	./$(PROJECT_NAME).exe
