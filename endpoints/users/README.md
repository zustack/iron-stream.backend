# Endpoints for users

{"exp":1726153719,"isAdmin":true,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYxNTM3MTksImlhdCI6MTcyMzU2MTcxOSwibmJmIjoxNzIzNTYxNzE5LCJzdWIiOjJ9.B1M7lPBPLbzC4D1qS4a1HZJjwPFuWY99vtBS_YhTo_o","userId":2}%

## Get admin user 
	isActiveParam := c.Query("a", "")
  isAdminParam := c.Query("admin", "")
  specialAppsParam := c.Query("special", "")
  verifiedParam := c.Query("verified", "")


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

## Create a new user
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "some",
    "password": "caro",
    "email": "carofricke@some.me",
    "name": "caro",
    "surname": "fricke",
    "pc": "some-pc",
    "os": "Mac"
  }'
```
## Login
```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"email": "carofricke@some.me", "password": "caro", "pc": "some-pc"}'
```
