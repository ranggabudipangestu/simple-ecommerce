export_env:
	export DB_USER="root" \
	export DB_PASS="root" \
	export DB_HOST="127.0.0.1" \
	export DB_PORT="3306" \
	export DB_NAME="simple-ecommerce"

test:
	go test ./internal/app/... -coverprofile=coverage.out -cover && go tool cover -func=coverage.out