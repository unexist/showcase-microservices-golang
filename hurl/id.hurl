GET http://localhost:8081/id
HTTP 201
[Asserts]
header "Content-Type" contains "utf-8"
jsonpath "$.id" == "string"
