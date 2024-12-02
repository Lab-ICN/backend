POST http://localhost:8080/api/v1/tokens
{
    "token": ""
}
HTTP 200
[Captures]
refreshToken: jsonpath "$['refreshToken']"
accessToken: jsonpath "$['accessToken']"

PUT http://localhost:8080/api/v1/tokens/self
{
    "refreshToken": "{{refreshToken}}"
}
HTTP 200

DELETE http://localhost:8080/api/v1/tokens/self
Authorization: Bearer {{accessToken}}
HTTP 200
