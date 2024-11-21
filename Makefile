run:
	@go run ./cmd/bookshop

tidy:
	@go mod tidy

lint:
	@golangci-lint run -c ./.golangci.yaml ./...

test:
	@go test ./... --cover --coverprofile=coverage.out

test/report: test
	@go tool cover -html=coverage.out

