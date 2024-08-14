# Endpoints for apps

## Create
```bash
# this endpoint need to be authenticated as admin
curl -X POST "http://localhost:8081/apps/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4" \
  -d '{"name": "Neovim", "process_name": "nvim", "os": "Mac", "is_active": true}'

```

## Get fist 50 apps && search by param && search by is_active
```bash
# this endpoint need to be authenticated as admin
# q= search query && (a=1 is_active=true, a=0 is_active=false)
curl -X GET "http://localhost:8081/apps/?cursor=0&q=&a=0" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
```

## Get app by id
```bash
# this endpoint need to be authenticated as admin
curl -X GET "http://localhost:8081/apps/get/1" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" | jq
```

## Update app by id
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/apps/update/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \
  -d '{"id": 1, "name": "Quick Time Player edited", "process_name": "QuickTime edited", "os": "Mac edited", "is_active": true}'
```

## Delete app by id
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/apps/update/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \
  -d '{"id": 1, "name": "Quick Time Player edited", "process_name": "QuickTime edited", "os": "Mac edited", "is_active": true}'
```

## Get apps by os
```bash
# this endpoint need to be authenticated as normal user
# this endpoint should return only the name and process_name
curl -X GET "http://localhost:8081/apps/normal-user/Mac" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" \ | jq
```

## Delete app by id
```bash
# this endpoint need to be authenticated as admin user
curl -X DELETE "http://localhost:8081/apps/delete/3" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1MjUwMjYsImlhdCI6MTcyMTkzMzAyNiwibmJmIjoxNzIxOTMzMDI2LCJzdWIiOjF9.qhWUBsobBK0TIWX2OD08HqlCas833r3bsQPKZTjlmU0" 
```
