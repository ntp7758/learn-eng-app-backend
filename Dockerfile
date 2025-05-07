FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .
RUN chmod +x wait-for-postgres.sh

RUN go build -o main ./cmd

# CMD ["/app/main"]