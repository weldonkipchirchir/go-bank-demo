# Use the official Golang image as base
FROM golang:1.22.1

# Set necessary environment variables
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate.linux-amd64 /usr/bin/migrate

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Run the Go app
CMD ["./main"]

# Expose the port on which your Go application listens
EXPOSE 8000
