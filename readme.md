
RUN Migration
migrate -path database/migrations -database "mysql://$DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/$DB_NAME" up