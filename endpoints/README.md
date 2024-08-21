# new endpoints test

- create user
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "agustfricke@gmail.com",
    "password": "admin",
    "name": "admin",
    "surname": "admin",
    "pc": "admin-pc",
    "os": "Linux"
  }'
```

- login 
```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"email": "agustfricke@gmail.com", "password": "admin", "pc": "admin-pc"}'
```

admin
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4NTg1MDQsImlhdCI6MTcyNDI2NjUwNCwibmJmIjoxNzI0MjY2NTA0LCJzdWIiOjF9.U4hlDaZeHAwNNHhR1KxC22GBP0mZ5D0IJr0GXWhJcBo

n
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4NTk5MjEsImlhdCI6MTcyNDI2NzkyMSwibmJmIjoxNzI0MjY3OTIxLCJzdWIiOjJ9.oZs_-gwZBnpVa3N5DyDXHpyurYhTpERVMOqj6kcvCVM

- create course
```bash
curl -X POST http://localhost:8081/courses/create \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4NTg1MDQsImlhdCI6MTcyNDI2NjUwNCwibmJmIjoxNzI0MjY2NTA0LCJzdWIiOjF9.U4hlDaZeHAwNNHhR1KxC22GBP0mZ5D0IJr0GXWhJcBo" \
  -H "Content-Type: multipart/form-data" \
  -F "title=test" \
  -F "description=Description test for Data Structures" \
  -F "author=agustfricke" \
  -F "duration=4 hours, 20 minutes" \
  -F "is_active=false" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "preview_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/test.mp4"
```
- get admin courses
```bash
curl -X GET "http://localhost:8081/courses/admin?q=&a="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4NTg1MDQsImlhdCI6MTcyNDI2NjUwNCwibmJmIjoxNzI0MjY2NTA0LCJzdWIiOjF9.U4hlDaZeHAwNNHhR1KxC22GBP0mZ5D0IJr0GXWhJcBo" | jq
```

- add course to user 
```bash
curl -X POST "http://localhost:8081/user/courses/1/1"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4NTg1MDQsImlhdCI6MTcyNDI2NjUwNCwibmJmIjoxNzI0MjY2NTA0LCJzdWIiOjF9.U4hlDaZeHAwNNHhR1KxC22GBP0mZ5D0IJr0GXWhJcBo"
```

-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" | jq

- get user courses 
```bash
curl -X GET "http://localhost:8081/user/courses?q="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4NTk5MjEsImlhdCI6MTcyNDI2NzkyMSwibmJmIjoxNzI0MjY3OTIxLCJzdWIiOjJ9.oZs_-gwZBnpVa3N5DyDXHpyurYhTpERVMOqj6kcvCVM" | jq
```

- desactivar todos los cursos a todos los usuarios(delete all the user_courses records)
```bash
curl -X DELETE "http://localhost:8081/user/courses/all"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" 
```

- desactivar curso especifico a todos los usuarios que lo contienen
```bash
curl -X DELETE "http://localhost:8081/user/courses/solo/1"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMzU0OTgsImlhdCI6MTcyMjU0MzQ5OCwibmJmIjoxNzIyNTQzNDk4LCJzdWIiOjF9.V1BbfsZ3-ZbNxJrU-TvrYrWmaWmsY128NHQYAZXV_Vc" 
```

