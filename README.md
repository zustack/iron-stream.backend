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

##### Want to run just a single test?
```bash
# for example run the login test
go test -v ./tests/users_test.go -run TestLogin
```

## Want to make curl requets?
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123",
    "email": "johndoe@example.com",
    "name": "John",
    "surname": "Doe",
    "pc": "My Computer",
    "os": "Windows 10"
  }'
```

```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "johndoe", "password": "securepassword123", "pc": "My Computer"}'
```
