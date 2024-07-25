# Iron Stream API

## Installation & setup
```bash
git clone https://github.com/zustack/iron-stream.backend.git ~/
cd ~/iron-stream.backend
```

##### cp the .env.example to .env and then edit it to your needs
```bash
cp ~/iron-stream.backend/.env.example .env
```

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
##### Register
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
##### Login
```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "admin", "pc": "admin-pc"}'
```

{"token":""}%

### Apps endpoints
##### Create
```bash
# this endpoint need to be authenticated as admin
curl -X POST "http://localhost:8081/apps/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \
  -d '{"name": "Quick Time Player", "process_name": "QuickTime", "os": "Mac", "is_active": true}'
```

##### Get fist 50 apps
```bash
# this endpoint need to be authenticated as admin
# cursor=0 means the first 50, cursor=1 means the second 50
# q= means search by name, process_name or os 
# the a=0 is not active, the a=1 is active
curl -X GET "http://localhost:8081/apps/?cursor=0&q=&a=" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
```

##### Search by is active
```bash
# this endpoint need to be authenticated as admin
# the a=0 is not active, the a=1 is active
curl -X GET "http://localhost:8081/apps/?cursor=0&q=&a=0" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
curl -X GET "http://localhost:8081/apps/?cursor=0&q=&a=1" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
```

##### Search by name, process_name, os or is_active
```bash
# this endpoint need to be authenticated as admin
curl -X GET "http://localhost:8081/apps/?cursor=0&q=Q&a=1" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
```

##### Get app by id
```bash
# this endpoint need to be authenticated as admin
curl -X GET "http://localhost:8081/apps/get/1" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
```

##### Update app by id
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/apps/update/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \
  -d '{"id": 1, "name": "Quick Time Player edited", "process_name": "QuickTime edited", "os": "Mac edited", "is_active": true}'
```

##### Delete app by id
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/apps/update/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \
  -d '{"id": 1, "name": "Quick Time Player edited", "process_name": "QuickTime edited", "os": "Mac edited", "is_active": true}'
```

##### Get apps by os
```bash
# this endpoint need to be authenticated as normal user
# this endpoint should return just the name and process_name
curl -X GET "http://localhost:8081/apps/normal-user/Mac" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \ | jq
```

##### Delete app by id
```bash
# this endpoint need to be authenticated as admin user
curl -X DELETE "http://localhost:8081/apps/delete/3" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" 
```
