# Build stage
# Use the official Golang image as base
FROM golang:1.22.1-alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main main.go

RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz 

# Run stage
FROM alpine
 
WORKDIR /app

COPY  --from=builder /app/main .

COPY --from=builder /app/migrate ./migrate 

COPY app.env .

COPY start.sh .

COPY wait-for.sh .

COPY db/migration ./migration

# Expose the port on which your Go application listens
EXPOSE 8000 

# Run the Go app
CMD ["/app/main"] 
ENTRYPOINT ["/app/start.sh"]