FROM golang:1.26.0-alpine3.22 AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates git

RUN go install github.com/swaggo/swag/cmd/swag@latest


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN swag init -g cmd/main/main.go

RUN GOOS=linux GOARCH=amd64 go build -o ./.bin/main ./cmd/main/main.go

FROM alpine:3.20 AS runner

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/.bin/main /

COPY --from=builder /app/docs .

COPY --from=builder /app/boot.yaml /app

COPY config config

COPY .env.prod .env.prod

ENV ENV_FILE=.env.prod

CMD ["/main"]
