GET http://localhost:3000/api/health
HTTP 200
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.status" == "RUNNING"
jsonpath "$.healthy" == true
jsonpath "$.operationId" exists
