FROM golang:1.26-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .

RUN ls -la && ls -la internal/

RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]