#Dockerfile
#Step 1 Build
FROM golang:1.11-alpine3.8 AS build
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /go/src/app
COPY . .
RUN go get github.com/gin-contrib/cors
RUN go get github.com/gin-gonic/gin
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/test ./dl.go
#Step 2 Final
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 23062
ENTRYPOINT /go/bin/test --port 23061