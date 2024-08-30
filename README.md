# Iron Stream Backend API

## Overview
Iron Stream REST API that powers the ui 
<a href="https://github.com/zustack/ui-iron-stream">Iron Stream</a>.
This project uses the <a href="https://go.dev">Go</a> programming language along with the 
web framework <a href="https://gofiber.io">Fiber</a>.

Iron Stream is currently in beta and under active development.

## Table of contents
- [Setup local development](#Setup)
- [API endpoints](#Endpoints) 
- [Tests](#Tests)
- [Deploy](#Deploy)

# Setup
To run this project locally, you will need: <a href="https://go.dev">Go</a>, 
<a href="https://www.ffmpeg.org">Ffmpeg</a> and <a href="https://www.sqlite.org">Sqlite</a>.

For install Go, you can download it from the [Go website](https://go.dev/dl/).

For Debian based systems, you can install Ffmpeg and Sqlite with the following commands:
```bash
sudo apt install ffmpeg sqlite3 -y
```

Clone the repository
```bash
git clone https://github.com/zustack/iron-stream.backend.git ~/
cd ~/iron-stream.backend
```

Create the directories for static files
```bash
mkdir web/uploads
cd web/uploads
mkdir files previews thumbnails tmp videos
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
## Users endpoints

### Signup request
```bash
curl -X POST "http://localhost:8081/users/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "agustfricke@gmail.com",
    "name": "Agustin",
    "surname": "Fricke",
    "password": "some-password",
    "pc": "agust@ubuntu",
    "os": "Linux"
  }'
```
This will send a confirmation email with a token used to verify the account in
the endpoint [verify](##verify)
### Signup response
```bash
200 OK
```

### Login request
```bash
curl -X POST "http://localhost:8081/users/login" \
     -H "Content-Type: application/json" \
     -d '{"email": "agustfricke@gmail.com", "password": "new-password", "pc": "agust@ubuntu"}'
```
### Login response
```json
{"exp":1727628111,"fullName":"Agustin Fricke","isAdmin":false,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2MjgxMTEsImlhdCI6MTcyNTAzNjExMSwibmJmIjoxNzI1MDM2MTExLCJzdWIiOjF9.ZBCbxsEbMoQS5legRGu1QArw3vcZV0jjqJ_f0u9l-0I","userId":1}
```

### Verify account request
```bash
curl -X POST "http://localhost:8081/users/verify" \
     -H "Content-Type: application/json" \
     -d '{"email": "agustfricke@gmail.com", "email_token": 902250}'
```
### Verify account response
```json
{"exp":1727627599,"fullName":"Agustin Fricke","isAdmin":false,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjc1OTksImlhdCI6MTcyNTAzNTU5OSwibmJmIjoxNzI1MDM1NTk5LCJzdWIiOjF9.dbpz5t6noMEW264uHL1AlbcOiVSrhbfiPvh9PwL1oSM","userId":1}
```

### Resend email request
This will send a new email with a token and update email_token field to the
user with the email provided.
```bash
curl -X POST "http://localhost:8081/users/resend/email/token/agustfricke@gmail.com" 
```
### Resend email response
```json
200 OK
```

### Delete account by email request
The user with the id 1 can't be deleted.
```bash
curl -X DELETE "http://localhost:8081/users/delete/account/by/email/agustfricke@protonmail.com"
```
### Delete account by email response
```json
200 OK
```

### Update password request
```bash
curl -X PUT "http://localhost:8081/users/update/password" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2MjgxMTEsImlhdCI6MTcyNTAzNjExMSwibmJmIjoxNzI1MDM2MTExLCJzdWIiOjF9.ZBCbxsEbMoQS5legRGu1QArw3vcZV0jjqJ_f0u9l-0I" \
     -d '{"password": "new-password"}'
```
### Update password response
```json
200 OK
```

### Get users request
This endpoint needs a valid admin token. The variables are:
- cursor: the cursor 
- q: the search query
- a: the active status
- admin: the admin status
- special: the special apps status
- verified: the verified status
```bash
curl -X GET "http://localhost:8081/users/admin?cursor=0&q=&a=&admin=&special=&verified="  \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjg3NjQsImlhdCI6MTcyNTAzNjc2NCwibmJmIjoxNzI1MDM2NzY0LCJzdWIiOjF9.8ptujSrdymZ7z5GdQAydWsVn4fmkgVg8rPOW4L37dxU" | jq
```
### Get users response
```json
{
  "data": [
    {
      "id": 1,
      "password": "$2a$10$OXFhJOHNkiLLc0EEmoasCOMH4FeXEqh4biody4Jo2K9WjDxTyfjMW",
      "email": "agustfricke@gmail.com",
      "name": "Agustin",
      "surname": "Fricke",
      "is_admin": true,
      "special_apps": false,
      "is_active": true,
      "email_token": 615068,
      "verified": true,
      "pc": "agust@ubuntu",
      "os": "Linux",
      "created_at": "30/08/2024 13:33:03"
    }
  ],
  "totalCount": 1,
  "previousId": null,
  "nextId": null
}
```

### Update user active status request
This endpoint needs a valid admin token. The user with the id 1 can't be updated.
```bash
curl -X PUT "http://localhost:8081/users/update/active/status/3" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjg3NjQsImlhdCI6MTcyNTAzNjc2NCwibmJmIjoxNzI1MDM2NzY0LCJzdWIiOjF9.8ptujSrdymZ7z5GdQAydWsVn4fmkgVg8rPOW4L37dxU" 
```
### Update user active status response
```json
200 OK
```

### Update all users active status request
This endpoint needs a valid admin token. The user with the id 1 can't be updated.
```bash
curl -X PUT "http://localhost:8081/users/update/all/active/status/true" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjg3NjQsImlhdCI6MTcyNTAzNjc2NCwibmJmIjoxNzI1MDM2NzY0LCJzdWIiOjF9.8ptujSrdymZ7z5GdQAydWsVn4fmkgVg8rPOW4L37dxU" 
```
### Update all users active status response
```json
200 OK
```

### Update special apps by user id request
This endpoint needs a valid admin token. 
```bash
curl -X PUT "http://localhost:8081/users/update/special/apps/user/1/true" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjg3NjQsImlhdCI6MTcyNTAzNjc2NCwibmJmIjoxNzI1MDM2NzY0LCJzdWIiOjF9.8ptujSrdymZ7z5GdQAydWsVn4fmkgVg8rPOW4L37dxU" 
```
### Update special apps by user id response
```json
200 OK
```

### Update admin status request
This endpoint needs a valid admin token. The user with the id 1 can't be updated.
```bash
curl -X PUT "http://localhost:8081/users/update/admin/status/3/true" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjg3NjQsImlhdCI6MTcyNTAzNjc2NCwibmJmIjoxNzI1MDM2NzY0LCJzdWIiOjF9.8ptujSrdymZ7z5GdQAydWsVn4fmkgVg8rPOW4L37dxU" 
```
### Update admin status response
```json
200 OK
```

# Tests
## Users tests
```bash
go test ./internal/database/users_test.go -v
go test ./api/inputs/users_test.go -v 
go test ./api/handlers/users_test.go -v
```

## Run all the available tests 
```bash
make test
```

# Deploy
```bash
make build
```

