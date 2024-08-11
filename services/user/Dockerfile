FROM golang:1.22-alpine

ENV PORT=8080
ENV DB_PATH=/app/sqlite.db

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary cmd/*.go

ENTRYPOINT ["/app/binary"]