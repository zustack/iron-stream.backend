# Endpoints for users

## Create a new user
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin",
    "email": "admin@admin.com",
    "name": "admin",
    "surname": "admin",
    "pc": "some-pc",
    "os": "Linux"
  }'
```
## Login
```bash
curl -X POST "http://localhost:8081/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "admin", "pc": "admin-pc"}'
```