# Contributing to Calidum Rotae Backend

Want to contribute to Calidum Rotae Backend? Here's an in-depth guide on how to do that.

## Dependencies

* [Go](https://go.dev/doc/install)
* [docker](https://docs.docker.com/get-docker/)
* [docker-compose](https://docs.docker.com/compose/install/)

### protoc

You'll also want to install the latest version of [protoc](https://grpc.io/docs/protoc-installation/) to be able to generate Go files from your Protobuf spec.

### other dependencies

Simply run this command to install the other dependencies:
```bash
$ make
```

## Generating gRPC Go files

To generate the gRPC Go files from the Protobuf spec, run this command:
```bash
$ make grpc
```

## Local development

To run the project locally using `docker-compose`, simply run this command:
```bash
$ make docker-compose
```

The API will then be available on port 3000.

## Testing the API
To test the API, run this command:
```bash
curl -X POST \
  http://localhost:3000/ \
  -H 'X-API-KEY: 1234' \
  -H 'Content-Type: application/json' \
  -d '{
  "Sender": {
    "FirstName": "firtname",
    "LastName": "lastname",
    "Email": "email@email.com",
    "PhoneNumber": "111-111-1111"
  },
  "RequestService": "service",
  "RequestDetails": "details"
}'
```
