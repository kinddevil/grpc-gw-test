
env GOOS=linux GOARCH=amd64 go build -o docker/app main.go \
  && cp -rf resources docker/ \
  && docker build -t grpc_template:latest docker \
  && rm -f docker/app && rm -rf docker/resources