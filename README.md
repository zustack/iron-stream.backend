# Iron Stream Backend API

## Overview
Iron Stream REST API that powers the ui 
<a href="https://github.com/zustack/ui-iron-stream">Iron Stream</a>

This project uses the <a href="https://go.dev">Go</a> programming language along with the 
web framework <a href="https://gofiber.io">Fiber</a>

Iron Stream is currently in beta and under active development.

[Setup local development](#Setup)
[API endpoints](#Endpoints)
[Tests](#Tests)
[Deploy](#Deploy)

# Setup
To run this project locally, you will need: <a href="https://go.dev">Go</a>, 
<a href="https://www.ffmpeg.org">Ffmpeg</a> and <a href="https://www.sqlite.org">Sqlite</a>.

For Debian based systems, you can install Ffmpeg and Sqlite with the following commands:
```bash
sudo apt install ffmpeg sqlite3 -y
```

Clone the repository
```bash
git clone https://github.com/zustack/iron-stream.backend.git ~/
cd ~/iron-stream.backend
```

For the environment variables, run this commands with the corresponding information.
```bash
export DB_DEV_PATH=/path/to/sqlite.db
export DB_TEST_PATH=/path/to/test_sqlite.db
export SECRET_KEY=someradomstring
export EMAIL_SECRET_KEY=emailsecret
export ROOT_PATH=/path/to/iron-stream.backend
```

Database
```bash
sqlite3 sqlite.db ".read tables.sql"
```

Then just execute:
```bash
make run
```

# Endpoints
[users](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/users) <br>
[apps](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/apps) <br>
[courses](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/courses) <br>
[history](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/history) <br>
[files](https://github.com/zustack/iron-stream.backend/tree/main/endpoints/files) <br>

# Tests
## User tests
```bash
go test ./internal/database/users_test.go -v
go test ./api/inputs/users_test.go -v 
go test ./api/handlers/users_test.go -v
```

## Run all the available tests 
```bash
go test ./... -v
```

# Deploy
```bash
make build
```

