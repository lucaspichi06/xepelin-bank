FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN  go mod tidy -compat=1.17
RUN go mod download
RUN go build -o ./xepelin-bank ./cmd/server/main.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/xepelin-bank .
ENTRYPOINT ["./xepelin-bank"]