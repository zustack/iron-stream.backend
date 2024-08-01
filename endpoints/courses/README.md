# Endpoints for courses

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik"}%
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik

## Upload a large file
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/chunk/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "chunkNumber=0" \
  -F "totalChunks=1" \
  -F "file=@/home/agust/Videos/test_edit.mp4"
```

## Response for Upload large file
```json
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/1cd35233-c3c5-4d6f-b155-1c1a74a1a133/test.mkv"}%
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/cc557fd0-239b-48f1-8f0c-a6d297bfca9c/test.mp4"}%
```

## Create course
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/create \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures" \
  -F "description=Description for Data Structures" \
  -F "author=agustfricke" \
  -F "duration=4 hours, 20 minutes" \
  -F "is_active=true" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "preview_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/cc557fd0-239b-48f1-8f0c-a6d297bfca9c/test.mp4"
```

## Get courses
```bash
curl -X GET "http://localhost:8081/courses/?cursor=0&q="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" | jq
```

## Get admin courses
```bash
# this endpoint need to be authenticated as admin
curl -X GET "http://localhost:8081/courses/admin?cursor=0&q=&a="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" | jq
```

## Add course to user
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/courses/add/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \
  -d '{"user_id": 1, "course_id": 1}'
```

## Update course
First you need to upload the video to http://localhost:8081/courses/chunk/upload and get 
the response that looks like this
```json
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/cc557fd0-239b-48f1-8f0c-a6d297bfca9c/test.mp4"}%
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/762716dc-f053-4150-86c0-e6678c0e263a/test_edit.mp4"}%
```

```bash
# this endpoint need to be authenticated as admin
curl -X PUT http://localhost:8081/courses/update \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "id=1" \
  -F "sort_order=4" \
  -F "title=Data Structures updated" \
  -F "description=Description for Data Structures updated" \
  -F "author=agustfricke updated" \
  -F "duration=6 hours, 09 minutes" \
  -F "is_active=false" \
  -F "thumbnail=@/home/agust/Pictures/edit.png" \
  -F "previewTmpDir=/home/agust/work/iron-stream/backend/web/uploads/tmp/762716dc-f053-4150-86c0-e6678c0e263a/test_edit.mp4"
```

## Delete course by id
```bash
# this endpoint need to be authenticated as admin user
curl -X DELETE "http://localhost:8081/courses/delete/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" 
```
