#!/bin/bash

go fmt ./...

docker-compose -f docker-compose.test.yml -f docker-compose.test.secrets.yml up --build
