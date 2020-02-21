#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8888 &
./server -port=8887 &
./server -port=8886 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait