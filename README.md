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
- from: from date example: 12/09/2024%00:00:00
- to: to date example: 12:09:2024%00:00:00

```bash
curl -X GET "http://localhost:8081/users/admin?cursor=0&q=&a=&admin=&special=&verified=&from=&to="  \
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

### Get current user request
```bash
curl -X GET "http://localhost:8081/users/current" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg1NzM4NjIsImlhdCI6MTcyNTk4MTg2MiwibmJmIjoxNzI1OTgxODYyLCJzdWIiOjF9.g9MhXLOzQIoSMPmfey4XlJvbOdknNysLLNINgblVcGU" | jq
```
### Get current user response
```json
{
  "created_at": "09/09/2024 16:46:56",
  "email": "agustfricke@gmail.com",
  "id": 1,
  "name": "Agustin",
  "surname": "Fricke"
}
```


<h1 id="courses-endpoints">Courses endpoints</h1>

### ChunkUpload request
```bash
curl -X POST http://localhost:8081/courses/chunk/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk0NDUxNzIsImlhdCI6MTcyNjg1MzE3MiwibmJmIjoxNzI2ODUzMTcyLCJzdWIiOjF9.rk1EW27xTcCXur3Vjh5_kYrGGtd-4D5_9e8icXz0ZaQ" \
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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures" \
  -F "description=Description for Data Structures" \
  -F "author=agustfricke" \
  -F "duration=4 hours, 20 minutes" \
  -F "is_active=true" \
  -F "price=420" \
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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
  -H "Content-Type: multipart/form-data" \
  -F "id=1" \
  -F "title=Data Structures edit" \
  -F "description=Description for Data Structures edit" \
  -F "author=agustfricke edit" \
  -F "duration=4 hours, 20 minutes edit" \
  -F "is_active=true" \
  -F "price=69" \
  -F "thumbnail=" \
  -F "old_thumbnail=/web/uploads/thumbnails/cd7af076-5f37-490c-bc34-998f4468b178.png" \
  -F "preview_tmp=" \
  -F "old_preview=/web/uploads/previews/511481d9-af34-40d7-8c1d-73ef42dfdb90/master.m3u8" \
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

### Get course by ID request
```bash
curl "http://localhost:8081/courses/solo/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \ | jq
```
### Get course by ID response
```bash
204 No Content
```

### Update course sort request
```bash
curl -X POST "http://localhost:8081/courses/sort" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" \
  -H "Content-Type: application/json" \
  -d '{
    "sort_courses": [
      {"id": 2, "sort_order": "8"},
      {"id": 3, "sort_order": "7"},
      {"id": 4, "sort_order": "6"},
      {"id": 5, "sort_order": "5"}
    ]
  }'
```
### Update course sort response
```bash
200 OK
```

### Update course active status request
```bash
curl -X PUT "http://localhost:8081/courses/update/active/5/true" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" 
```
### Update course active status response
```json
200 OK
```

### admin courses
This endpoint needs a valid admin token. The variables are:
- cursor: the cursor 
- q: the search query
- a: the active status("": all, "1": active, "0": non-active)
```bash
curl -X GET "http://localhost:8081/courses/admin?cursor=0&q=&a="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" | jq
```
### admin courses response
```json
[
  {
    "id": 2,
    "title": "Data Structures",
    "description": "Description for Data Structures",
    "author": "agustfricke",
    "thumbnail": "/web/uploads/thumbnails/ed73325d-9f37-4b94-b6f1-b98f258d4211.png",
    "preview": "/home/agust/work/iron-stream/backend/web/uploads/previews/627e9d78-0698-4cf4-990c-afb0491276ed/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "4 hours, 20 minutes",
    "is_active": true,
    "sort_order": 8,
    "created_at": "31/08/2024 15:53:24"
  },
  {
    "id": 3,
    "title": "new course",
    "description": "some long ",
    "author": "Created by  me!",
    "thumbnail": "/web/uploads/thumbnails/fcf3d710-808b-4604-bc5a-08ce422fd65a.png",
    "preview": "/home/agust/work/iron-stream/backend/web/uploads/previews/1d88fd82-bc15-4cd1-8423-c4c79172c070/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "lalala",
    "is_active": true,
    "sort_order": 7,
    "created_at": "31/08/2024 15:55:49"
  },
  {
    "id": 7,
    "title": "some dsa",
    "description": "some dsa",
    "author": "some dsa",
    "thumbnail": "/web/uploads/thumbnails/f08b6aad-d407-4402-8912-3a9025a13ea0.png",
    "preview": "/web/uploads/previews/b2927bed-767e-47f8-b4af-9acca1ace537/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "hello",
    "is_active": true,
    "sort_order": 7,
    "created_at": "31/08/2024 16:09:44"
  },
  {
    "id": 4,
    "title": "foo",
    "description": "soso",
    "author": "Created by ",
    "thumbnail": "/web/uploads/thumbnails/c29bff9a-d440-40ff-b513-07bb8030e91e.png",
    "preview": "/web/uploads/previews/65b1e13b-0ab7-4f9b-aa07-12beb2febaa8/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "sosoo",
    "is_active": false,
    "sort_order": 6,
    "created_at": "31/08/2024 15:57:44"
  },
  {
    "id": 6,
    "title": "test",
    "description": "test",
    "author": "Created by ",
    "thumbnail": "/web/uploads/thumbnails/d7af27d6-5936-451a-a185-84320f55781b.jpg",
    "preview": "/web/uploads/previews/ea6a32ae-14cd-4065-be6d-2dd9a1069c90/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "sd",
    "is_active": true,
    "sort_order": 6,
    "created_at": "31/08/2024 16:05:13"
  },
  {
    "id": 5,
    "title": "some",
    "description": "somsom",
    "author": "Created by ",
    "thumbnail": "/web/uploads/thumbnails/24129485-951b-480d-90fb-3f5aae5b799c.png",
    "preview": "/web/uploads/previews/e2a97d9a-db66-49cb-9fda-7c07ae0e67e5/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "somso",
    "is_active": true,
    "sort_order": 5,
    "created_at": "31/08/2024 16:02:10"
  }
]
```

### Add course to user request
```bash
curl -X POST "http://localhost:8081/user/courses/8/6" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw"
```
### Add course to user response
```json
200 OK
```

### Get user courses request
```bash
curl -X GET "http://localhost:8081/user/courses?q=" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" | jq
```
### Get user courses response 
```json
[
  {
    "id": 2,
    "title": "Data Structures",
    "description": "Description for Data Structures",
    "author": "agustfricke",
    "thumbnail": "/web/uploads/thumbnails/ed73325d-9f37-4b94-b6f1-b98f258d4211.png",
    "preview": "/home/agust/work/iron-stream/backend/web/uploads/previews/627e9d78-0698-4cf4-990c-afb0491276ed/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "4 hours, 20 minutes",
    "is_active": true,
    "sort_order": 8,
    "created_at": "31/08/2024 15:53:24",
    "is_user_enrolled": false
  },
  {
    "id": 3,
    "title": "new course",
    "description": "some long ",
    "author": "Created by  me!",
    "thumbnail": "/web/uploads/thumbnails/fcf3d710-808b-4604-bc5a-08ce422fd65a.png",
    "preview": "/home/agust/work/iron-stream/backend/web/uploads/previews/1d88fd82-bc15-4cd1-8423-c4c79172c070/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "lalala",
    "is_active": true,
    "sort_order": 7,
    "created_at": "31/08/2024 15:55:49",
    "is_user_enrolled": false
  },
  {
    "id": 4,
    "title": "foo",
    "description": "soso",
    "author": "Created by ",
    "thumbnail": "/web/uploads/thumbnails/c29bff9a-d440-40ff-b513-07bb8030e91e.png",
    "preview": "/web/uploads/previews/65b1e13b-0ab7-4f9b-aa07-12beb2febaa8/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "sosoo",
    "is_active": false,
    "sort_order": 6,
    "created_at": "31/08/2024 15:57:44",
    "is_user_enrolled": false
  },
  {
    "id": 5,
    "title": "some",
    "description": "somsom",
    "author": "Created by ",
    "thumbnail": "/web/uploads/thumbnails/24129485-951b-480d-90fb-3f5aae5b799c.png",
    "preview": "/web/uploads/previews/e2a97d9a-db66-49cb-9fda-7c07ae0e67e5/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "somso",
    "is_active": true,
    "sort_order": 5,
    "created_at": "31/08/2024 16:02:10",
    "is_user_enrolled": false
  },
  {
    "id": 6,
    "title": "test",
    "description": "test",
    "author": "Created by ",
    "thumbnail": "/web/uploads/thumbnails/d7af27d6-5936-451a-a185-84320f55781b.jpg",
    "preview": "/web/uploads/previews/ea6a32ae-14cd-4065-be6d-2dd9a1069c90/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "sd",
    "is_active": true,
    "sort_order": 6,
    "created_at": "31/08/2024 16:05:13",
    "is_user_enrolled": false
  },
  {
    "id": 7,
    "title": "some dsa",
    "description": "some dsa",
    "author": "some dsa",
    "thumbnail": "/web/uploads/thumbnails/f08b6aad-d407-4402-8912-3a9025a13ea0.png",
    "preview": "/web/uploads/previews/b2927bed-767e-47f8-b4af-9acca1ace537/master.m3u8",
    "rating": 0,
    "num_reviews": 0,
    "duration": "hello",
    "is_active": true,
    "sort_order": 7,
    "created_at": "31/08/2024 16:09:44",
    "is_user_enrolled": false
  }
]
```

### Delete all user_courses request
This will delete all the user_courses
```bash
curl -X DELETE "http://localhost:8081/user/courses/all" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw"
```
### Delete all user_courses response
```json
200 OK 
```

### Delete user_courses by courseId request
This will "deactivate" the specific course to all the users.
```bash
curl -X DELETE "http://localhost:8081/user/courses/solo/7" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw"
```
### Delete user_courses by courseId response
```json
200 OK 
```

### Delete user_courses by courseId & userId request
This will "deactivate" the specific course to the specific user.
```bash
curl -X DELETE "http://localhost:8081/user/courses/solo/8/7" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw"
```
### Delete user_courses by courseId response
```json
200 OK 
```

## Videos

### Create video request
```bash
curl -X POST http://localhost:8081/videos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures 420" \
  -F "description=Description for Data Structures 420" \
  -F "course_id=1" \
  -F "duration=4 hours, 20 minutes" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "video_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/test.mp4"
```
### Create video response
```json
200 OK
```

### Get admin videos by course id request
```bash
curl -X GET "http://localhost:8081/videos/admin/1?q=" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" | jq
```
### Get admin videos response
```json
[
  {
    "id": 2,
    "title": "Data Structures 420",
    "description": "Description for Data Structures 420",
    "video_hls": "/web/uploads/previews/f059bce1-a9bc-479c-a738-f1d87ec63618/master.m3u8",
    "thumbnail": "/web/uploads/thumbnails/0c364957-c0d9-4451-93e9-53e95c06897d.png",
    "length": "119",
    "duration": "4 hours, 20 minutes",
    "views": 0,
    "course_id": "1",
    "created_at": "02/09/2024 14:54:44",
    "video_resume": ""
  }
]
```

### Get video feed request
```bash
curl -X GET "http://localhost:8081/videos/feed/1?q=" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884"
```
### Get video feed response
```json
[
  {
    "id": 2,
    "title": "Data Structures 420",
    "description": "Description for Data Structures 420",
    "video_hls": "/web/uploads/previews/f059bce1-a9bc-479c-a738-f1d87ec63618/master.m3u8",
    "thumbnail": "/web/uploads/thumbnails/0c364957-c0d9-4451-93e9-53e95c06897d.png",
    "length": "119",
    "duration": "4 hours, 20 minutes",
    "views": 0,
    "course_id": "1",
    "created_at": "",
    "video_resume": ""
  }
]
```

### Delete video
```bash
curl -X DELETE "http://localhost:8081/videos/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884"
```
### Delete video response
```json
204 No Content
```

### Update video request
```bash
curl -X PUT http://localhost:8081/videos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
  -H "Content-Type: multipart/form-data" \
  -F "title=edit edit Data Structures" \
  -F "description=edit edit Description for Data Structures" \
  -F "id=4" \
  -F "duration=edit edit 420" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "old_thumbnail=/web/uploads/thumbnails/b62d3919-0864-454c-bee3-e7c959c5b0a3.png" \
  -F "old_video_hls=/web/uploads/videos/1/5b0f2475-6d88-4a75-bab1-284bdd375d29/master.m3u8" \
  -F "video_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/test.mp4"
```
### Update video response
```json
200 OK
```

### Get current video request
```bash
curl -X GET "http://localhost:8081/history/current/3" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" | jq
```
### Get current video response
```json
{
  "history_id": "",
  "resume": "",
  "video": {
    "id": 2,
    "title": "Data Structures 420",
    "description": "Description for Data Structures 420",
    "video_hls": "/web/uploads/previews/f059bce1-a9bc-479c-a738-f1d87ec63618/master.m3u8",
    "thumbnail": "/web/uploads/thumbnails/0c364957-c0d9-4451-93e9-53e95c06897d.png",
    "length": "4 hours, 20 minutes",
    "duration": "119",
    "views": 0,
    "course_id": "1",
    "created_at": "02/09/2024 14:54:44",
    "video_resume": ""
  }
}
```

### Watch new video request
```bash
curl -X PUT "http://localhost:8081/history/watch" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
     -d '{"history_id": "21", "video_id": 7, "course_id": "3", "resume": "69", "current_video_id": 7}'
```
### Watch new video response
```json
{"id":22,"video_id":"","course_id":0,"user_id":0,"video_resume":"187.68","created_at":""}%
```


## History
### Get user history request
```bash
curl -X GET "http://localhost:8081/history" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg1NzM4NjIsImlhdCI6MTcyNTk4MTg2MiwibmJmIjoxNzI1OTgxODYyLCJzdWIiOjF9.g9MhXLOzQIoSMPmfey4XlJvbOdknNysLLNINgblVcGU" | jq
```

### Get user history response
```json
[
  {
    "history_id": 1,
    "video_id": 1,
    "course_id": 1,
    "user_id": 1,
    "video_resume": "120.8820645161291",
    "history_date": "10/09/2024 15:36:42",
    "video_title": "video 1",
    "description": "some long ass",
    "video_hls": "/web/uploads/videos/1/8e0a486f-ba98-4847-92e1-cb677d40bb8d/master.m3u8",
    "thumbnail": "/web/uploads/thumbnails/dd046c69-07b3-4192-a756-0a1bfd5edf11.png",
    "duration": "4",
    "length": "187",
    "views": 3,
    "video_created": "10/09/2024 15:35:33"
  },
  {
    "history_id": 2,
    "video_id": 2,
    "course_id": 1,
    "user_id": 1,
    "video_resume": "",
    "history_date": "10/09/2024 15:38:18",
    "video_title": "2222",
    "description": "sdsdsd",
    "video_hls": "/web/uploads/videos/1/fed969d9-bc0a-445f-a327-6c76fa6a0c62/master.m3u8",
    "thumbnail": "/web/uploads/thumbnails/3416c96d-6dd1-4673-9ea6-5e4c3950dd6e.png",
    "duration": "ddf",
    "length": "119",
    "views": 1,
    "video_created": "10/09/2024 15:36:28"
  }
]
```


### Update history
```bash
curl -X PUT "http://localhost:8081/history/update" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
     -d '{"id": "23", "resume": "6969"}'
```
### Update history response
```json
200 OK
```

## Apps
### Create app request
```bash
curl -X POST "http://localhost:8081/apps/create" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" \
     -d '{"name": "OBS", "process_name": "obs", "is_active": true, "execute_always": true}'
```
```json
200 OK
```

### Delete app request
```bash
curl -X DELETE "http://localhost:8081/apps/delete/6" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" 
```
### Delete app response
```json
204 No Content
```

### Update app status request
```bash
curl -X PUT "http://localhost:8081/apps/update" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" \
     -d '{"id": 6, "name": "obs", "process_name": "obs", "is_active": true, "execute_always": true}'
```
### Update app status response
```json
200 OK
```

### Get admin apps request
```bash
curl -X GET "http://localhost:8081/apps/admin?q=&a=" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" | jq
```
### Get admin apps response
```json
[
  {
    "id": 4,
    "name": "OBS",
    "process_name": "obs.app",
    "is_active": true,
    "execute_always": false,
    "created_at": "03/09/2024 14:24:22"
  },
  {
    "id": 1,
    "name": "FireFox",
    "process_name": "firefox.app",
    "is_active": true,
    "execute_always": false,
    "created_at": "03/09/2024 13:31:09"
  }
]
```

### Get forbidden apps request
```bash
curl -X GET "http://localhost:8081/apps/forbidden" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" | jq
```
### Get forbidden apps response
```json
[
  {
    "id": 0,
    "name": "FireFox",
    "process_name": "firefox.app",
    "is_active": false,
    "execute_always": false,
    "created_at": ""
  },
  {
    "id": 0,
    "name": "OBS",
    "process_name": "obs.app",
    "is_active": false,
    "execute_always": false,
    "created_at": ""
  }
]
```

### Update app status request
```bash
curl -X PUT "http://localhost:8081/apps/update/status/5/false" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" 
```

### Update app status response
```json
200 OK
```

### Update app ea status request
```bash
curl -X PUT "http://localhost:8081/apps/update/ea/7/false" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" 
```

### Update app ea status response
```json
200 OK
```

## User apps
### Create user app
```bash
curl -X POST "http://localhost:8081/user/apps/create/1/5" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" 
```
### Get user apps(admin to get all the apps with field related)
```bash
curl -X GET "http://localhost:8081/user/apps/user/apps/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" | jq
```
### Get user apps response
```json
[
  {
    "id": 1,
    "name": "FireFox",
    "process_name": "firefox.app",
    "is_active": false,
    "execute_always": false,
    "created_at": "",
    "is_user_enrolled": false
  },
  {
    "id": 4,
    "name": "OBS",
    "process_name": "obs.app",
    "is_active": false,
    "execute_always": false,
    "created_at": "",
    "is_user_enrolled": false
  },
  {
    "id": 5,
    "name": "ZSH",
    "process_name": "zsh.app",
    "is_active": false,
    "execute_always": false,
    "created_at": "",
    "is_user_enrolled": true
  }
]
```

### Delete user app request
```bash
curl -X DELETE "http://localhost:8081/user/apps/delete/user/apps/1/5" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM"
```
### Delete user app response
```json
204 No Content
```

## Files
### Create file
```bash
curl -X POST http://localhost:8081/files \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDEzNzIsImlhdCI6MTcyNTQwOTM3MiwibmJmIjoxNzI1NDA5MzcyLCJzdWIiOjF9.sjhoW7xJcSsBcC4XqVYDyqGyas74OGwLz7dPobmmIjc" \
  -H "Content-Type: multipart/form-data" \
  -F "videoID=1" \
  -F "page=1" \
  -F "path=@/home/agust/Pictures/test.png"
```
### Create video response
```json
200 OK
```

### Get files request
```bash
curl -X GET "http://localhost:8081/files/2/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk" | jq
```
### Get files response
```json
{
  "files": [
    {
      "id": 1,
      "path": "/web/uploads/files/43fb91b7-ecd6-4e0e-a43a-7ff0c55753a7.svg",
      "video_id": "2",
      "page": "1",
      "created_at": "03/09/2024 22:13:59"
    },
    {
      "id": 2,
      "path": "/web/uploads/files/4ff3de3f-a344-4bd5-ad55-9fabd992b0d4.svg",
      "video_id": "2",
      "page": "1",
      "created_at": "03/09/2024 22:27:57"
    },
    {
      "id": 3,
      "path": "/web/uploads/files/2d883b7d-cae3-4e74-9aa8-5267725d55a1.svg",
      "video_id": "2",
      "page": "1",
      "created_at": "03/09/2024 22:28:20"
    }
  ],
  "pageCount": 1
}
```
### Delete file request
```bash
curl -X DELETE "http://localhost:8081/files/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc5NzMwNDAsImlhdCI6MTcyNTM4MTA0MCwibmJmIjoxNzI1MzgxMDQwLCJzdWIiOjF9.e4RGUkZS75k8KkCNpczLhtl7mCtpMBz0EVXWH2fdZSM" | jq
```
### Delete file response
```json
204 No Content
```

## Reviews
### Create review request
```bash
curl -X POST "http://localhost:8081/reviews" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwNTU1NzksImlhdCI6MTcyNTQ2MzU3OSwibmJmIjoxNzI1NDYzNTc5LCJzdWIiOjF9.E0Q8S1gl7Ka560mM-8mSE0sepJQZHhbXIry9Qc7JMHA" \
  -d '{
    "course_id": "2",
    "description": "I like this course v2",
    "rating": 4.3
  }'
```
### Create review response
```json
200 OK
```

### Admin reviews
```bash
curl -X GET "http://localhost:8081/reviews/admin?q=&p=" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk" | jq
```
### Admin reviews response
```json
[
  {
    "id": 2,
    "course_id": 2,
    "user_id": 1,
    "author": "Agustin Fricke",
    "description": "I like this course v2",
    "rating": 4.3,
    "public": false,
    "course_title": "Data Structures",
    "created_at": "04/09/2024 12:30:44"
  }
]
```

### Reviews public by course id
```bash
curl -X GET "http://localhost:8081/reviews/public/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk" | jq
```

```json
[
  {
    "id": 2,
    "course_id": 2,
    "user_id": 1,
    "author": "Agustin Fricke",
    "description": "I like this course v2",
    "rating": 4.3,
    "public": true,
    "course_title": "Data Structures",
    "created_at": "04/09/2024 12:30:44"
  }
]
```

### Update public status  request
```bash
curl -X PUT "http://localhost:8081/reviews/update/public/1/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk"
```
### Update public status response
```json
200 OK
```

### Delete Review request
```bash
curl -X DELETE "http://localhost:8081/reviews/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk"
```
### Delete Review response
```json
204 No Content
```

## Notes
### Create note request
```bash
curl -X POST "http://localhost:8081/notes/2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwNTU1NzksImlhdCI6MTcyNTQ2MzU3OSwibmJmIjoxNzI1NDYzNTc5LCJzdWIiOjF9.E0Q8S1gl7Ka560mM-8mSE0sepJQZHhbXIry9Qc7JMHA" \
  -d '{
    "body": "Note for some video",
    "video_title": "The video title",
    "time": "420.69"
  }'
```
### Create note response
```json
200 OK
```

### Get notes response
```bash
curl -X GET "http://localhost:8081/notes/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk" | jq
```
### Get notes response
```json
[
  {
    "id": 1,
    "body": "Note for some video",
    "video_title": "The video title",
    "time": "420.69",
    "course_id": "1",
    "user_id": 1
  }
]
```
### Edit note request
```bash
curl -X PUT "http://localhost:8081/notes/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwNTU1NzksImlhdCI6MTcyNTQ2MzU3OSwibmJmIjoxNzI1NDYzNTc5LCJzdWIiOjF9.E0Q8S1gl7Ka560mM-8mSE0sepJQZHhbXIry9Qc7JMHA" \
  -d '{
    "body": "edit Note for some video"
  }'
```
### Edit note response
```json
200 OK
```

### Delete note response
```bash
curl -X DELETE "http://localhost:8081/notes/1" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk" 
```
### Delete note response
```json
204 No Content
```

## Notifications
### Get admin notifications
```bash
curl -X GET "http://localhost:8081/notifications" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMDQ5NTMsImlhdCI6MTcyNTQxMjk1MywibmJmIjoxNzI1NDEyOTUzLCJzdWIiOjF9._VLiflVTJ5tTP2Li0l0XE7TPFfqONkB341m6F_XFgOk"  | jq
```

## Policy
### Create policy request
```bash
curl -X POST "http://localhost:8081/policy" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg1MDMyNDksImlhdCI6MTcyNTkxMTI0OSwibmJmIjoxNzI1OTExMjQ5LCJzdWIiOjF9.0nxdOjcn6yyxLzkVwsF0sLKP3FiXfZQaTFd7QPnJ7mk" \
  -d '{
    "content": "La actual **Política** de Privacidad establece los límites y parámetros que Carolina Fricke emplea para el uso y protección del conjunto de datos de carácter personal proporcionados por los usuarios de “https://www.carofricke.com”, comprometida en la seguridad de todos ellos, para evitar su uso indebido.",
    "p_type": "text"
  }'
```
### Create policy response
```json
200 OK
```

### Get policy request
```bash
curl -X GET "http://localhost:8081/policy" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg1MDMyNDksImlhdCI6MTcyNTkxMTI0OSwibmJmIjoxNzI1OTExMjQ5LCJzdWIiOjF9.0nxdOjcn6yyxLzkVwsF0sLKP3FiXfZQaTFd7QPnJ7mk" | jq
```
### Get policy response
```json
[
  {
    "id": 1,
    "content": "Información básica sobre protección de datos.",
    "p_type": "li",
    "created_at": "09/09/2024 16:47:54"
  }
]
```

### Delete policy request
```bash
curl -X DELETE "http://localhost:8081/policy/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg1MDMyNDksImlhdCI6MTcyNTkxMTI0OSwibmJmIjoxNzI1OTExMjQ5LCJzdWIiOjF9.0nxdOjcn6yyxLzkVwsF0sLKP3FiXfZQaTFd7QPnJ7mk"
```
### Delete policy response
```json
204 No Content
```

## User log
### Create logout user log request
```bash
curl -X POST "http://localhost:8081/log/user/logout" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM"
```
### Create logout user log response
```json
200 OK
```

### Create login user log request
- This is created at login

### Get user log request
```bash
curl -X GET "http://localhost:8081/log/user/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM" | jq
```
### Get user log response
```json
[
  {
    "id": 5,
    "content": "The user has logged in.",
    "l_type": "1",
    "user_id": 1,
    "created_at": "11/09/2024 16:42:40"
  },
  {
    "id": 6,
    "content": "The user has logged out.",
    "l_type": "1",
    "user_id": 1,
    "created_at": "11/09/2024 16:42:47"
  },
  {
    "id": 7,
    "content": "The app OBS was open while watching the video Data structures.",
    "l_type": "3",
    "user_id": 1,
    "created_at": "11/09/2024 16:42:57"
  },
  {
    "id": 8,
    "content": "The apps OBS, Google Chrome where open while watching the video Data structures.",
    "l_type": "3",
    "user_id": 1,
    "created_at": "11/09/2024 16:43:16"
  }
]
```

### Create log found apps request
```bash
curl -X POST "http://localhost:8081/log/user/found/apps" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM" \
  -d '{
      "video_title": "Data structures",
      "apps": [
      { "name": "OBS" },
      { "name": "Google Chrome" }
      ]
  }'
```
### Create log found apps response
```json
200 OK
```

## Admin log
### Create admin log request
```bash
curl -X POST "http://localhost:8081/log/admin" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM" \
  -d '{
      "content": "The user pepe@gmail.com was deleted.",
      "l_type": "3"
  }'
```
### Create admin log response
```json
204 No Content
```
### Get Admin log request
```bash
curl -X GET "http://localhost:8081/log/admin" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM" | jq
```
### Get Admin log response
```json
[
  {
    "id": 1,
    "content": "The user pepe@gmail.com was deleted.",
    "l_type": "3",
    "created_at": "12/09/2024 15:46:03"
  }
]
```

## Statistics

### Get user statistics request
- Get basic data from September 2024
```bash
curl -X GET "http://localhost:8081/users/stats?from=2024-09-09&to=2024-09-16"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM" | jq
```
### Get user statistics response
```json
[
  {
    "date": "2024/09/09",
    "windows": 0,
    "mac": 0,
    "linux": 1,
    "all": 1
  },
  {
    "date": "2024/09/10",
    "windows": 0,
    "mac": 0,
    "linux": 0,
    "all": 0
  },
  {
    "date": "2024/09/11",
    "windows": 0,
    "mac": 0,
    "linux": 0,
    "all": 0
  },
  {
    "date": "2024/09/12",
    "windows": 0,
    "mac": 0,
    "linux": 0,
    "all": 0
  },
  {
    "date": "2024/09/13",
    "windows": 0,
    "mac": 0,
    "linux": 0,
    "all": 0
  },
  {
    "date": "2024/09/14",
    "windows": 0,
    "mac": 0,
    "linux": 0,
    "all": 0
  },
  {
    "date": "2024/09/15",
    "windows": 0,
    "mac": 0,
    "linux": 0,
    "all": 0
  },
  {
    "date": "2024/09/16",
    "windows": 0,
    "mac": 0,
    "linux": 1,
    "all": 1
  }
]
```

### Get course profit statistics request
```bash
curl -X GET "http://localhost:8081/courses/stats?from=19/09/2024&to=19/09/2024"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjg2NzU0NDQsImlhdCI6MTcyNjA4MzQ0NCwibmJmIjoxNzI2MDgzNDQ0LCJzdWIiOjF9.UoA5afV7b4blotaaFfcyKnojjVMWupBRzAwF08dBPbM" | jq
```

### Get course profit statistics response
```json
[
  {
    "Title": "new course",
    "Profit": 138
  }
]
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

## Files
```bash
cwebp image.jpg -o image.webp
cwebp some-png.png -o image-png.webp
cwebp -q 50 imagen.jpg -o imagen.webp
du -h *
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
