POST http://localhost:8080/todo
{
    "description": "string",
    "title": "string"
}
HTTP 201
[Asserts]
header "Content-Type" contains "utf-8"
jsonpath "$.title" == "string"
jsonpath "$.description" == "string"

GET http://localhost:8080/todo
HTTP 200
