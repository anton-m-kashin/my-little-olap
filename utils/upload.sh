#!/bin/sh

HOST="127.0.0.1:8080"
ENDPOINT="upload"

curl \
    --silent \
    --request POST \
    --data @mock-data.json \
    "http://${HOST}/${ENDPOINT}"
