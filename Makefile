run:
	@go run ./cmd/guestbook

lint:
	@golangci-lint run -c ./.golangci.yaml ./...

test:
	@go test ./... --cover --coverprofile=coverage.out

test/report: test
	@go tool cover -html=coverage.out

