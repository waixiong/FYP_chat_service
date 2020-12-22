# Where you should clone this project
At `$GOPATH/src/rj`, git clone this project.

# Pre
- Download Golang
- Download Docker
- Download protobuf compiler (http://google.github.io/proto-lens/installing-protoc.html)

# For TLS cert
https://workaround.org/ispmail/jessie/create-certificate
### In development
`openssl req -newkey rsa:4096 -nodes -sha512 -x509 -days 3650 -nodes -out ./key/certs/mycert.pem -keyout ./key/private/mykey.pem`

### In production
`openssl genrsa -out ./key/private/mykey.pem 4096`

`openssl req -new -key ./key/private/mykey.pem -out ./key/certs/mycert.csr`


# Docker (Currently for AWS)
Run this to pull
`docker pull amazon/dynamodb-local`

Run this to start
`docker run -t -p 10001:8000 amazon/dynamodb-local`


# AWS dynamo
https://docs.aws.amazon.com/en_pv/sdk-for-go/v1/developer-guide/using-dynamodb-with-go-sdk.html


# Server
### Some Dependecies
`go get github.com/aws/aws-sdk-go/aws`

`go get google.golang.org/grpc`

`go get google.golang.org/api/oauth2/v2`


### Folder Structure
- api
  - proto
    - ${services}.proto
- cmd
  - <server>
    - main.go (setup server)
  - key (dir that save key)
- pkg
  - api
    - ${services}
      - ${services}.pb.go
  - service
    - ${services}
      - ${services}.go (file, main part of the service)
      - src (dir for any extra go program)


For more of google services, please look at https://github.com/googleapis/google-api-go-client 
