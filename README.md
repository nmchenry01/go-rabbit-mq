# go-rabbit-mq

## Overview

The purpose of this repository is to POC interacting with RabbitMQ using Golang as well as sending/unmarshalling XML based payloads. 

## Getting started

1. Navigate to the root of the directory and run `docker-compose up --build`

2. Navigate to the `consumer/` directory and run `go run main.go`

3. Navigate to the `producer/` directory and run `go run main.go`

You should see some messages be send from the producer application to the consumer application. Feel free to play around with this repo and add more functionality as needed.