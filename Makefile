MOCKERY_BIN := $(GOPATH)/bin/mockery

.PHONY: serve tidy test mock

serve:
	go run cmd/api/main.go
tidy:
	go mod tidy && go mod download && go mod vendor
test:
	go run cmd/test/main.go
mock:
	@echo "Generating mocks for interface $(interface) in directory $(dir)..."
	@$(MOCKERY_BIN) --name=$(interface) --dir=$(dir) --output=./internal/mocks
	cd ./internal/mocks && \
	mv $(interface).go $(filename).go
mig-up:
	go run cmd/migration/main.go -up
mig-down:
	go run cmd/migration/main.go -down
coverage:
	go test -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
seed:
	go run cmd/seed/main.go
swag:
	swag init -g main.go -dir ./cmd/api,./internal/controller/handler,./internal/controller/repository/user,./internal/controller/repository/address,./internal/controller/dto/response,./internal/controller/dto/request,./internal/controller/service,pkg/utils,./internal/core,./internal/helper,./internal/route