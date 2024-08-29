# User related endpoints

## Signup
```bash
curl -X POST http://localhost:8081/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "agustfricke@protonmail.com",
    "name": "Agustin",
    "surname": "Fricke",
    "password": "some-password",
    "pc": "agust@ubuntu",
    "os": "Linux"
  }'
```
UPDATE users SET is_admin = true WHERE id = 1

## Signup response
<h1>Login</h1>
<h5>Login resquest</h5>
```bash
curl -X POST "http://localhost:8081/users/login" \
     -H "Content-Type: application/json" \
     -d '{"email": "agustfricke@gmail.com", "password": "some-password", "pc": "some-pc"}'
```
<h5>Login response</h5>
```json
{"exp":1727549719,"fullName":"Agustin Fricke","isAdmin":true,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc1NDk3MTksImlhdCI6MTcyNDk1NzcxOSwibmJmIjoxNzI0OTU3NzE5LCJzdWIiOjF9.ISV92LzwLtw8msrRCc8iI5i32f207No6qlsN6YziL3M","userId":1}%
```


```bash
curl -X PUT "http://localhost:8081/courses/add/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4" \
  -d '{"user_id": 3, "course_id": 2}'
  ```

{"exp":1726161802,"isAdmin":true,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4","userId":4}%

```bash
curl -X GET "http://localhost:8081/users/admin?cursor=0&q=&a=&admin=&special=&verified="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4" | jq
```

  app.Put("deactivate/all/courses", middleware.AdminUser, handlers.DeactivateAllCoursesForAllUsers)
  app.Put("deactivate/course/for/all/users/:id", middleware.AdminUser, handlers.DeactivateSpecificCourseForAllUsers)

```bash
curl -X PUT "http://localhost:8081/deactivate/all/courses"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4" 
```

```bash
curl -X PUT "http://localhost:8081/deactivate/course/for/all/users/1"  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4" 
```

```bash
curl -X PUT "http://localhost:8081/update/active/status"  \
  -H "Content-Type: multipart/form-data" \
  -F "isActive=true" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNjE4MDIsImlhdCI6MTcyMzU2OTgwMiwibmJmIjoxNzIzNTY5ODAyLCJzdWIiOjR9.cn0fUJUF6ZFE6Iklxt-CL1KR2_uJ5eHWfX4iOFQdKi4" 
```

## Login
```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"email": "carofricke@some.me", "password": "caro", "pc": "some-pc"}'
```
