# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

ARG version

# Add Maintainer Info
LABEL maintainer="Alec Scott <alecbcs@github.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -ldflags "-s -w -X github.com/arken/arken/config.Version=$version" -o arken .

# Start again with minimal envoirnment.
FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/arken /app/arken

# Command to run the executable
CMD ["/app/arken", "daemon"]
