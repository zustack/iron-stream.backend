# Endpoints for files

## Create file
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/files \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "videoID=2" \
  -F "page=1" \
  -F "path=@/home/agust/Pictures/test.png"
```

## Get files by page and video id
```bash
# this endpoint need to be authenticated as noraml user
curl -X GET "http://localhost:8081/files?page=1&videoID=2" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
```

## Delete file by id
```bash
# this endpoint need to be authenticated as noraml user
curl -X DELETE "http://localhost:8081/files?id=1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
```
