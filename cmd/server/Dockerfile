FROM golang:1.17.1-alpine AS builder
WORKDIR /go/src
COPY . .
WORKDIR /go/src/cmd/server
RUN go get && go mod tidy
RUN go build -o ../../bin/server
WORKDIR /go/src
CMD ["bin/server"]

FROM alpine:latest
WORKDIR root
COPY --from=builder /go/src/bin ./bin
EXPOSE 8080
CMD ["bin/server"]
