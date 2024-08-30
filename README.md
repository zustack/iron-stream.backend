# Iron Stream Backend API

## User tests
```bash
go test ./internal/database/users_test.go -v
go test ./api/inputs/users_test.go -v 
go test ./api/handlers/users_test.go -v
```

## Run all the tests available
```bash
go test ./... -v
```

## Overview
This project uses the Go programming language along with the web framework Fiber 
and uses FFmpeg under the hood to convert and stream audio and video.

The frontend of this project can be found at [zustack/ui-iron-stream](https://github.com/zustack/ui-iron-stream.

The project is under active development so there may be some functionality missing or broken.

## Requirements
- Go 
- Ffmpeg
- Sqlite

## Installation & setup
```bash
git clone https://github.com/zustack/iron-stream.backend.git ~/
cd ~/iron-stream.backend
```

## Environment variables
For the environment variables, run this commands with the corresponding information.
```bash
export DB_DEV_PATH=/path/to/sqlite.db
export DB_TEST_PATH=/path/to/test_sqlite.db
export SECRET_KEY=someradomstring
export EMAIL_SECRET_KEY=emailsecret
export ROOT_PATH=/path/to/iron-stream.backend
```

## Database
#### Dev database 
```bash
# dev
sqlite3 sqlite.db ".read tables.sql"
```
#### Test database 
```bash
# test
sqlite3 test_sqlite.db ".read tables.sql"
```

## Make file instuctions
```bash
make run
```
```bash
make build
```
```bash
make test
```

## Endpoints
[users](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/users) <br>
[apps](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/apps) <br>
[courses](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/courses) <br>
[history](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/history) <br>
[files](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/files) <br>
