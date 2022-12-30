#! /bin/bash

COUNT=1000

make_requests() {
  cookie='Cookie: _ym_isad2=5'
  rand_str=$(openssl rand -hex 20)

  response=$(
    curl --location --request POST 'http://localhost:8080' \
      --header "$cookie" \
      --header 'Content-Type: text/plain' \
      --data-raw "$rand_str"
  )
  echo "$response"

  curl --location --request GET "$response" \
    --header "$cookie"

  echo
  curl --location --request GET 'http://localhost:8080/ping' \
    --header "$cookie"

  echo
  response=$(
    curl --location --request GET 'http://localhost:8080/api/user/urls' \
      --header "$cookie"
  )
  echo "$response"
}

for ((i = 1; i <= COUNT; i++)); do
  make_requests
done
