# Videos endpoints 

## Upload a large file
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/chunk/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "chunkNumber=0" \
  -F "totalChunks=1" \
  -F "file=@/home/agust/Videos/test.mp4"
```

## Response for Upload large file
```json
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/efb005be-ebc8-4bcf-9209-e9e7292e2c3f/test.mp4"}%
```

## Create video
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/videos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures 420" \
  -F "description=Description for Data Structures 420" \
  -F "course_id=1" \
  -F "length=4 hours, 20 minutes" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "video_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/c628f54d-b00b-4bc2-a86b-d9fc3e260d60/test.mp4"
```

## get videos by course id
```bash
curl -X GET "http://localhost:8081/videos/1/?cursor=0&q="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" | jq
```

## Update video
```bash
# this endpoint need to be authenticated as admin
curl -X PUT http://localhost:8081/videos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures update" \
  -F "description=Description for Data Structures update" \
  -F "id=1" \
  -F "length=4 hours, 20 minutes update" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "video_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/c628f54d-b00b-4bc2-a86b-d9fc3e260d60/test.mp4"
```

## delete video by course id
```bash
curl -X DELETE "http://localhost:8081/videos/1"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" | jq
```

## update video view and create history

```bash
curl -X PUT "http://localhost:8081/videos/views" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" \
  -d '{"video_id": 2, "course_id": 1, "video_resume": "069420"}'
```
