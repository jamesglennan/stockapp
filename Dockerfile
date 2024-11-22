# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# set go mod in its own step for caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main .

# Use multi-step build to keep the image small/secure
FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/main /

CMD ["/main"]