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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc4OTEyOTgsImlhdCI6MTcyNTI5OTI5OCwibmJmIjoxNzI1Mjk5Mjk4LCJzdWIiOjF9.uYzFWle0Apbk89vQ3azD8pe5yBghw8EAx_Jx_p_h884" \
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

### Get course by ID request
```bash
curl "http://localhost:8081/courses/solo/2" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc3MTYzMjUsImlhdCI6MTcyNTEyNDMyNSwibmJmIjoxNzI1MTI0MzI1LCJzdWIiOjh9._j6dGt0wiBPizAn3dCYnr1NKAksUIi7SYQJ1xmoH_Fw" | jq
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
