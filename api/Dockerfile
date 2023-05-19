# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /api

# Copy files and download dependencies
ADD . /api
RUN go mod download
RUN go mod tidy

# Build executable
RUN CGO_ENABLED=0 GOOS=linux go build -o /certsplace-api

# Expose Port
EXPOSE 8080

# Run
CMD ["/certsplace-api"]