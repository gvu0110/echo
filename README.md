# echo

```
curl localhost:8000

curl "localhost:8000/cats/string?name=Tom&age=2"
curl "localhost:8000/cats/json?name=Tom&age=2"
curl "localhost:8000/cats/abc?name=Tom&age=2"

curl -X POST -d '{"name":"Kitty","age":"3"}' localhost:8000/addcat
curl -X POST -d '{"name":"Pitbull","age":"4"}' localhost:8000/adddog
curl -X POST -H 'Content-Type: application/json' -d '{"name":"Amber","age":"1"}' localhost:8000/addhamster

curl -X GET localhost:8000/admin/main
curl -X GET -u adam:12345 localhost:8000/admin/main

curl -v -X GET localhost:8000

curl -v -c cookie.txt -X GET "localhost:8000/login?username=adam&password=12345"
```