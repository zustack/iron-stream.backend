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
    - [Users](#Users-endpoints)
    - [Courses](#Courses-endpoints)
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
     -d '{"email": "agustfricke@gmail.com", "password": "some-password", "pc": "agust@ubuntu"}'
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


<h1 id="courses-endpoints">Courses endpoints</h1>

### ChunkUpload request
```bash
curl -X POST http://localhost:8081/courses/chunk/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" \
  -H "Content-Type: multipart/form-data" \
  -F "chunkNumber=0" \
  -F "totalChunks=1" \
  -F "uuid=42069" \
  -F "file=@/home/agust/Videos/test.mp4"
```

### ChunkUpload response
```bash
/home/agust/work/iron-stream/backend/web/uploads/tmp/42069/test.mp4
```

### Create course request
```bash
curl -X POST http://localhost:8081/courses/create \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures" \
  -F "description=Description for Data Structures" \
  -F "author=agustfricke" \
  -F "duration=4 hours, 20 minutes" \
  -F "is_active=true" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "preview_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/42069/test.mp4"
```

### Create course response
```bash
200 OK
```

### Update course request
```bash
curl -X PUT "http://localhost:8081/courses/update" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" \
  -H "Content-Type: multipart/form-data" \
  -F "id=1" \
  -F "title=Data Structures edit" \
  -F "description=Description for Data Structures edit" \
  -F "author=agustfricke edit" \
  -F "duration=4 hours, 20 minutes edit" \
  -F "is_active=true" \
  -F "thumbnail=" \
  -F "old_thumbnail=/web/uploads/thumbnails/2da28e22-33b6-42ec-9e21-8a46023ecb2b.png" \
  -F "preview_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/42069/test.mp4" \
  -F "old_preview=/web/uploads/tmp/42069/test.png" \
```
### Update course response
```bash
200 OK
```

### Delete course by ID request
```bash
curl -X DELETE "http://localhost:8081/courses/delete/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" 
```
### Delete course by ID response
```bash
204 No Content
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
### Run the tests(recommended)
If you want to run the tests before deploying, you can execute:
```bash
make test
```

### Edit server port
Edit in `cmd/main.go` to set the port to 80.
```go
package main

import (
	"iron-stream/api"
	"iron-stream/internal/database"
	"log"
)

func main() {
	database.ConnectDB("DB_DEV_PATH")
	app := api.Setup()
	log.Fatal(app.Listen(":80")) // this here!
}
```

### Build the app
Build the app, this make command will create a binary file called iron-stream and
create the tables in the database.
```bash
make deploy
```

### Setup systemd service
Create a new file called `iron-stream.service` with the the corresponding information.
```bash
[Unit]
Description=iron-stream

[Service]
Environment=DB_DEV_PATH=/path/to/sqlite.db 
Environment=SECRET_KEY=someradomstring
Environment=EMAIL_SECRET_KEY=emailsecret
Environment=ROOT_PATH=/path/to/iron-stream.backend
Type=simple
Restart=always
RestartSec=5s
ExecStart=/path/to/iron-stream/binary/iron-stream

[Install]
WantedBy=multi-user.target
```

Move the file to `/etc/systemd/system`.
```bash
sudo mv iron-stream.service /etc/systemd/system
```

Run the service and check the status.
```bash
sudo service iron-stream start
sudo service iron-stream status
```

You should see something like this:
```bash
● iron-stream.service - iron-stream
     Loaded: loaded (/lib/systemd/system/iron-stream.service; disabled; vendor preset: enabled)
     Active: active (running) since Fri 2024-08-30 14:44:31 -03; 15min ago
   Main PID: 45633 (iron-stream)
      Tasks: 6 (limit: 18713)
     Memory: 2.3M
        CPU: 849ms
     CGroup: /system.slice/iron-stream.service
             └─45633 /path/to/iron-stream/binary/iron-stream
```

### Enable the Firewall 
Enable the port 80
```bash
sudo ufw enable
sudo ufw allow 80/tcp
```
Check the status
```bash
sudo ufw status
```

You should see something like this:
```bash
Status: active

To                         Action      From
--                         ------      ----
80/tcp                     ALLOW       Anywhere
80/tcp (v6)                ALLOW       Anywhere (v6)

```

### SSL certificates
I personally let Cloudflare handle the SSL certificates. 
Just point your domain to your server's IP address, as it's running on HTTP port (80).

### Example DNS record:
The `Proxy Status: Redirected via Proxy` will add a SSL certificate.
- Type: A
- Name: domain.com
- Value: 420.69.420.73
- Proxy Status: Redirected via Proxy
