# Quick cURLs

Quick commands to test the API.

## User service

### Register a user

```bash
curl -X 'POST' 'http://localhost:8100/v1/register' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{
  "username": "aaa",
  "password": "password",
  "email": "aaa@gmail.com",
  "phone_number": "+62891234567891",
  "first_name": "Wendy",
  "last_name": "Surya Wijaya"
}' | jq
```

### Login

```bash
curl -X 'POST' 'http://localhost:8100/v1/login' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{
  "username": "wendy",
  "password": "password"
}' | jq
```

### Get user profile

```bash
curl -X 'GET' 'http://localhost:8100/v1/profile' -H 'accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <token>' | jq
```

### Update user profile

```bash
curl -X 'PUT' 'http://localhost:8100/v1/profile' -H 'accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <token>' -d '{
  "first_name": "Wendys"
}' | jq
```

### Get verification email

```bash
curl -X GET 'http://localhost:8100/v1/profile/email/verification' -H "Authorization: Bearer <token>" | jq
```

## OPA

### Check if a user has access to a resource

```bash
curl -X "POST" http://localhost:8181/v1/data/policies/example/policy/allow -d @input.json
```

## Movie search service

### Search movie

```bash
curl -sS -X "GET" "http://localhost:8103/v1/movie?sort_by=release_date_desc&limit=9&theater_id=1234567890&include_upcoming=true&genre=drama&keyword=godfather" | jq
```

## Search movie by cursor

```bash
curl -sS -X "GET" "http://localhost:8103/v1/movie?cursor=eyJzb3J0X2J5IjoicmVsZWFzZV9kYXRlX2Rlc2MiLCJsaW1pdCI6OSwidGhlYXRlcl9pZCI6IjEyMzQ1Njc4OTAiLCJpbmNsdWRlX3VwY29taW5nIjp0cnVlLCJnZW5yZSI6ImRyYW1hIiwia2V5d29yZCI6ImdvZGZhdGhlciIsImxhc3Rfc2Vlbl9zb3J0X3ZhbHVlIjoxMCwibGFzdF9zZWVuX2lkIjoiNjdhY2E2MDkxMWI3NjhjMzI1MDUyOTQwIn0" | jq
```

## Movie service

## Admin search movie

```bash
res=$(curl -s 'http://localhost:8100/v1/login' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{
    "username": "wendy",
    "password": "XAMSIFgjZrSXKQHbIGCNEtWTRqValloifOhlPfEjBIIkGSMeUQ"
}');
curl -v 'http://localhost:8103/v1/admin/movies?theater_id=1&sort_by=relevance&page=1&page_size=10&keyword=truman' -H "Authorization: Bearer $( echo $res | jq -r '.result.access_token' | tr -d '\"' )"
```

### Admin search movie + propagate W3C trace context

```bash
#!/bin/bash

# Call the login endpoint and store headers in a temporary file.
res=$(curl -s -D headers.txt 'http://localhost:8100/v1/login' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{
    "username": "wendy",
    "password": "XAMSIFgjZrSXKQHbIGCNEtWTRqValloifOhlPfEjBIIkGSMeUQ"
}');

# Extract the 'traceparent' header from the response headers.
traceparent=$(grep -i '^traceparent:' headers.txt | awk '{print $2}' | tr -d '\r')

if [ -z "$traceparent" ]; then
  echo "No traceparent header found in the login response."
  exit 1
fi

echo "Extracted traceparent: $traceparent"

# Use the extracted traceparent header in the search endpoint request.
search_response=$(curl -v 'http://localhost:8103/v1/admin/movies?theater_id=1&sort_by=relevance&page=1&page_size=10&keyword=truman' -H "Authorization: Bearer $( echo $res | jq -r '.result.access_token' | tr -d '\"' )" -H "traceparent: $traceparent")

echo "Search response: $search_response"

rm -f headers.txt
```
