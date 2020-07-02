#!/bin/bash

go run main.go > /dev/null || exit &
PID=$?

sleep 3

curl -XPOST http://localhost:1323/password -d 'password=password' && echo 
curl -XPOST http://localhost:1323/password -d 'password=n0onepwn3d' && echo 

kill $PID
