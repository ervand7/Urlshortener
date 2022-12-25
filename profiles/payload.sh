#! /bin/bash

COUNT=1000

make_requests() {
  rand_str=$(openssl rand -hex 20)
  response=$(
    curl --location --request POST 'http://localhost:8080' \
      --header 'Cookie: _ym_isad2=5' \
      --header 'Content-Type: text/plain' \
      --data-raw "$rand_str"
  )
  echo "$response"

  curl --location --request GET "$response" \
    --header 'Cookie: _ym_isad2=5'

  echo
  curl --location --request GET 'http://localhost:8080/ping' \
    --header 'Cookie: _ym_isad2=5'

  echo
  response=$(
    curl --location --request GET 'http://localhost:8080/api/user/urls' \
      --header 'Cookie: _ym_isad2=5'
  )
  echo "$response"
}

for ((i = 1; i <= COUNT; i++)); do
  make_requests
done
