# Login
POST http://localhost:8080/user/login
HTTP 200
[Asserts]
header "Content-Type" contains "utf-8"
jsonpath "$.token" matches /[0-9]+/
jsonpath "$.token" matches /[0-9a-z-]+/
[Captures]
bearer: jsonpath "$.token"

# Self
GET http://localhost:8080/user/self
Authorization: Bearer {{bearer}}
HTTP 200
[Asserts]
header "Content-Type" contains "utf-8"
jsonpath "$.name" == "Default User"
jsonpath "$.token" == {{bearer}}
