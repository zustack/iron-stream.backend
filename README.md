# Iron Stream API

## Installation & setup
```bash
git clone https://github.com/zustack/iron-stream.backend.git ~/
cd ~/iron-stream.backend
```

## Environment variables
For the environment variables, run this command with the corresponding information.
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

## Tests
```bash
# register test
go test -v ./tests/users_test.go -run TestRegister
```
```bash
# login test
go test -v ./tests/users_test.go -run TestLogin
```
```bash
# create app test
go test -v ./tests/apps_test.go -run TestCreateApp
```

## Endpoints
All the endpoints can be found on the **endpoints/*/README.md** dir.
[users](https://github.com/zustack/www-iron-stream/endpoints/users/README.md)
[apps](https://github.com/zustack/www-iron-stream/endpoints/apps/README.md)
[courses](https://github.com/zustack/www-iron-stream/endpoints/courses/README.md)
