#Dockerfile
#Step 1 Build
FROM golang:1.11-alpine3.8 AS build
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /go/src/app
COPY . .
COPY ./templates ./bin/templates
COPY ./static ./bin/static
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/static
RUN go get github.com/gin-contrib/cors
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/test ./ul.go
EXPOSE 23061
ENTRYPOINT /go/src/app/bin/test --port 23061