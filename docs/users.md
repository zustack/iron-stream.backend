# Users endpoints

The first thing would be to export the environment variable `API_URL`.
By default it run on port 8081 but feel free to change it.
```bash
export API_URL="http://localhost:8081"
```


## User Signup

### Permissions
Any user can access this resource.

### Endpoint
`POST ${API_URL}/users/signup`

### Request Example
To register, send a POST request using the following `curl` command:

```bash
curl -X POST "${API_URL}/users/signup" \
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

### Field Descriptions
- **email**: The user's email address.
- **name**: The user's first name.
- **surname**: The user's last name.
- **password**: The user's password.
- **pc**: The name of the user's computer.
- **os**: The user's operating system.

### Verification Process
After submitting the request, a confirmation email will be sent with a token used to verify the account. To verify your account, use the following endpoint: [Verify Account](#verify-account).

### Response
If the registration is successful, you will receive the following response:

```json
201 Created
```

---

## Verify Account

### Permissions
Any user can access this resource.

When the user [registers](#user-signup), an email with a verification code will be sent to them.

### Endpoint
`POST ${API_URL}/users/verify`

### Request Example
To verify your account, send a POST request using the following `curl` command:

```bash
curl -X POST "${API_URL}/users/verify" \
  -H "Content-Type: application/json" \
  -d '{
      "email": "agustfricke@gmail.com", 
      "email_token": 674488
  }'
```

### Response
If the verification is successful, you will receive a response like the following:

```json
{
    "exp": 1727627599,
    "fullName": "Agustin Fricke",
    "isAdmin": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjc1OTksImlhdCI6MTcyNTAzNTU5OSwibmJmIjoxNzI1MDM1NTk5LCJzdWIiOjF9.dbpz5t6noMEW264uHL1AlbcOiVSrhbfiPvh9PwL1oSM",
    "userId": 1
}
```

---


## Login request
`Permission`: Any user can access this resource.
By default, the user with id 1 should be the super administrator. 
To achieve this, you can execute the following command:
```sql
UPDATE users SET is_admin = true WHERE id = 1;
```
Make the login request:
```bash
curl -X POST "${API_URL}/users/login" \
     -H "Content-Type: application/json" \
     -d '{
        "email": "agustfricke@protonmail.com", 
        "password": "some-password", 
        "pc": "agust@ubuntu"
    }'
```
Response: 
```json
{
    "exp":1727628111,
    "fullName":"Agustin Fricke",
    "isAdmin":true,
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2MjgxMTEsImlhdCI6MTcyNTAzNjExMSwibmJmIjoxNzI1MDM2MTExLCJzdWIiOjF9.ZBCbxsEbMoQS5legRGu1QArw3vcZV0jjqJ_f0u9l-0I",
    "userId":1
}
```
Now you can export the token in the ACCESS_TOKEN variable. to hit the endpoints
that need authentication.
```bash
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzAxNDQ3NDcsImlhdCI6MTcyNzU1Mjc0NywibmJmIjoxNzI3NTUyNzQ3LCJzdWIiOjF9._vnlD2OJgr21jCHV7FbceSU1eBqwuLKfl6PG9U72vSQ"
```

## Verify account 
`Permission`: Any user can access this resource.
When the user [registers](##register), an email with a verification code was 
sent to them.
```bash
curl -X POST "${API_URL}/users/verify" \
     -H "Content-Type: application/json" \
     -d '{
         "email": "agustfricke@gmail.com", 
         "email_token": 674488
         }'
```
Response:
```json
{
    "exp":1727627599,
    "fullName":"Agustin Fricke",
    "isAdmin":true,
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc2Mjc1OTksImlhdCI6MTcyNTAzNTU5OSwibmJmIjoxNzI1MDM1NTk5LCJzdWIiOjF9.dbpz5t6noMEW264uHL1AlbcOiVSrhbfiPvh9PwL1oSM",
    "userId":1
}
```

## Resend email 
`Permission`: Any user can access this resource.
`Email`: the user email.
Will send a new code to verify the email.
```bash
curl -X POST "${API_URL}/users/resend/email/token/{email}" 
```
Response:
```json
200 OK
```

## Delete account by email 
`Info`: The user with the id 1 can't be deleted.
`Permission`: Only the admin and owner of the account can access this resource.
`Email`: the user email
```bash
curl -X DELETE "${API_URL}/users/delete/account/by/email/{email}" \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" 
```
Response:
```json
200 OK
```

## Update password 
`Permission`: Only the owner of the account can access this resource.
```bash
curl -X PUT "${API_URL}/users/update/password" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" \
     -d '{"password": "new-password"}'
```
Response:
```json
200 OK
```

## Get users 
`Permission`: Only admin user can access this resource.
`cursor`: determine how many rows to skip in the result set before 
returning the next batch of users.
Example:
If you want to fetch the first 10 users, you'd set limit = 10 and cursor = 0.
If you want to fetch the next 10 users, you'd set limit = 10 and cursor = 10, and so on.
`q`: the search query
Includes:
email, name, surname, os and created_at 
`a`: the active status
Can be true or false
`admin`: the admin status
Can be true or false
`special`: the special apps status
Can be true or false
`verified`: the verified status
Can be true or false
`from`: from date 
Example: from=28/09/202420%01:59:21
`to`: to date 
Example: to=30/09/2024%2001:59:21
```bash
curl -X GET "${API_URL}/users/admin?cursor=0&q=&a=&admin=&special=&verified=&from=28/09/202420%01:59:21&to=30/09/2024%2001:59:21" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" | jq
curl -X GET "${API_URL}/users/admin?cursor=0&q=&a=&admin=&special=&verified=&from=&to="  \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" | jq
```
Response:
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

## Update user active status 
`Permission`: Only admin user can access this resource.
`Info`: The user with the id 1 can't be updated.
`userId`: the user id
It will change the current active status of the user. 
if true it will update to false and vice versa.
```bash
curl -X PUT "${API_URL}/users/update/active/status/{userId}" \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" 
```
Response:
```json
200 OK
```

## Update all users active status 
`Permission`: Only admin user can access this resource.
`Info`: The user with the id 1 can't be updated.
`active`: can be true or false.
```bash
curl -X PUT "${API_URL}/users/update/all/active/status/{active}" \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" 
```
Response:
```json
200 OK
```

## Update special apps by user id 
`Permission`: Only admin user can access this resource.
`userId`: the user id
`active`: can be true or false.
The active param will update the special apps status of the user.
```bash
curl -X PUT "${API_URL}/users/update/special/apps/user/{userId}/{active}" \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" 
```
Response:
```json
200 OK
```

## Update admin status 
`Permission`: Only admin user can access this resource.
`Info`: The user with the id 1 can't be updated.
`userId`: the user id
`isAdmin`: can be true or false.
Will update the admin status of the user with id `userId`.
```bash
curl -X PUT "${API_URL}/users/update/admin/status/{userId}/{isAdmin}" \
     -H "Authorization: Bearer ${ACCESS_TOKEN}" 
```
Response:
```json
200 OK
```

## Get current user 
`Permission`: Only logged user can access this resource.
```bash
curl -X GET "${API_URL}/users/current" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" | jq
```
Response:
```json
{
  "created_at": "09/09/2024 16:46:56",
  "email": "agustfricke@gmail.com",
  "id": 1,
  "name": "Agustin",
  "surname": "Fricke"
}
```

