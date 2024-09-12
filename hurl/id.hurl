GET http://localhost:8081/id
HTTP 201
[Asserts]
header "Content-Type" contains "utf-8"
jsonpath "$.id" matches "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

