FROM golang:1.22-alpine

RUN apk add --no-cache gcc musl-dev


ENV PORT=8080
ENV DB_PATH=/app/sqlite.db
ENV CGO_ENABLED=1

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary cmd/*.go

ENTRYPOINT ["/app/binary"]