# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# Final stage
FROM golang:1.25-alpine
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]