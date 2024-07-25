# Iron Stream 

## Make file instuctions

##### Run for development
```bash
make run
```

##### Build
```bash
make build
```

##### Run every test with a new testing database
```bash
make test
```

##### Apply database migrations
```bash
# test
sqlite3 test_sqlite.db ".read tables.sql"
# dev
sqlite3 sqlite.db ".read tables.sql"
```

##### Want to run just a single test?
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

## Want to make curl requets?
- Register 
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin",
    "email": "admin@admin.com",
    "name": "admin",
    "surname": "admin",
    "pc": "some-pc",
    "os": "Linux"
  }'
```

- Login
```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "some-password", "pc": "admin-pc"}'
```

- Create new app
```bash
curl -X POST http://localhost:8081/apps/create \
  -H "Content-Type: application/json" \
  # put the jwt token here
  -H "Authorization: Bearer " \
  -d '{"name": "OBS", "test-process": "obs.exe", "os": "Windows", "is-active": true}'
```
