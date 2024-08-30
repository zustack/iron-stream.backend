run:
	go run ./cmd/main.go

deploy:
	go build -o iron-stream cmd/main.go
	sqlite3 sqlite.db ".read tables.sql"

test:
	go test ./... -v

