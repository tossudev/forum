# Use official Go Docker image
# from https://hub.docker.com/_/golang

# It is worth considering a lightweight variant
# such as the alpine image. (https://hub.docker.com/_/golang#image-variants)
# However, since this is labeled as highly experimental,
# the defacto golang image seems like a more stable option
FROM golang:1.27 AS build

# Create and cd working directory (good practice)
WORKDIR /app

# Acquire dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy all files from project root to Docker container
# Uses .dockerignore to ignore unnecessary files (such as *.db or *.md)
COPY . .

# Build program
RUN go build -tags fts5 -o forum ./cmd

# Expose HTTP port for the container
EXPOSE 8080

# Run server, reset (and seed) database
# TODO: Have options for pre-seeded database and completely empty database
CMD ["./forum", "--reset"]

