test:
	go test ./internal/app/... -coverprofile=coverage.out -cover && go tool cover -func=coverage.out

run:
	go build -o binary && ./binary

dev:
	go run main.go