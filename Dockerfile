# Use an official Golang runtime as the base image
FROM golang:1.17-alpine

# Install GCC and other necessary dependencies
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/web
ENV DB_PATH=/app/pkg/db/database.db
EXPOSE 8080
CMD ["./main"]
