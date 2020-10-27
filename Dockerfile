#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o hydration-app cmd/hydration-app/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/app/hydration-app /app
ENTRYPOINT ./app
LABEL Name=mongodbtest Version=0.0.1
EXPOSE 8080