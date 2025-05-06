FROM golang:1.24.2 AS auth-builder

WORKDIR /app

COPY . .

RUN cd cmd/auth && CGO_ENABLED=0 GOOS=linux go build -o auth-service main.go

FROM alpine:latest

WORKDIR /

COPY --from=auth-builder /app/cmd/auth/auth-service ./

COPY --from=auth-builder /app/internal/config/.env ./

ENTRYPOINT ["./auth-service"]