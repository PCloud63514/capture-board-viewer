PROJECT_NAME := capture-board-selector
PKG_PATH_BATCH := ./

verify:
	go mod tidy
	go mod vendor

ci:
	go mod tidy
	go mod vendor
	CGO_ENABLED=0 go build -o $(PROJECT_NAME).exe $(PKG_PATH_BATCH)

run:
	$(MAKE) ci
	./$(PROJECT_NAME).exe
