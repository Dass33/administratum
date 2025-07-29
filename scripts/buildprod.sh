#!/bin/bash

cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o administratum
