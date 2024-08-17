FROM golang:1.22-alpine

RUN apk add --no-cache gcc musl-dev


ENV AUTH_SERVICE_ADDRESS=auth-service:8081
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary cmd/*.go

ENTRYPOINT ["/app/binary"]