# install dependencies
go mod tidy

# deploy mongodb
docker-compose up

# config .env file
.env

# start the service
go run main.go

# api doc
https://apifox.com/apidoc/shared-12fd0fcc-7e42-455b-9cbb-021682aeeed4/api-106900727