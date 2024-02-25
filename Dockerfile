# syntax=docker/dockerfile:1

FROM golang:alpine3.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . /app

# Build
RUN go build -o forecast ./cmd/main.go

# Run
CMD ["./forecast"]