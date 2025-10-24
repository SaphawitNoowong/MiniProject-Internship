FROM golang:1.22-alpine AS builder

WORKDIR /app

# cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build static linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

# --- runtime image ---
FROM alpine:3.20

WORKDIR /app

# add CA certs and non-root user
RUN apk --no-cache add ca-certificates \
    && adduser -D -g '' appuser

COPY --from=builder /app/main /app/main

EXPOSE 5000
USER appuser

CMD ["./main"]
