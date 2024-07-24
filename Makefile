run:
	sqlite3 sqlite.db ".read tables.sql"
	go run ./cmd/main.go

build:
	go build cmd/main.go

test:
	sqlite3 test_sqlite.db ".read tables.sql"
	go test ./tests -v
