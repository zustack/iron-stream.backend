# Endpoints for courses

## Upload a large file
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/chunk/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \
  -H "Content-Type: multipart/form-data" \
  -F "chunkNumber=0" \
  -F "totalChunks=1" \
  -F "file=@/home/agust/Videos/test.mkv"
```

## Response for Upload large file
```json
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/1cd35233-c3c5-4d6f-b155-1c1a74a1a133/test.mkv"}%
```

## Create course
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/create \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures" \
  -F "description=Description for Data Structures" \
  -F "author=agustfricke" \
  -F "duration=4 hours, 20 minutes" \
  -F "is_active=true" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "preview_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/1cd35233-c3c5-4d6f-b155-1c1a74a1a133/test.mkv"
```

## Get courses
```bash
curl -X GET "http://localhost:8081/courses/?cursor=0&q="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \ | jq
```

## Get admin courses
```bash
# this endpoint need to be authenticated as admin
curl -X GET "http://localhost:8081/courses/admin?cursor=0&q=&a="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \ | jq
```

## Add course to user
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/courses/add/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \
  -d '{"user_id": 1, "course_id": 1}'
```
