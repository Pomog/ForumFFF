# Use the official Go image as a parent image
FROM golang:latest
LABEL maintainer="Your Name dverves; ypanasiu"
LABEL version="1.0"
LABEL description="This is a custom Docker image for my application."
# Set GOPATH to an empty string within the container
ARG GOPATH=
# Copy your Go application source code into the container
COPY ./ ./
# Build the Go application
RUN go build -o ffforum cmd/web/*.go
# Expose a port (if your application listens on a specific port)
EXPOSE 8080
# Command to run your Go application
CMD ["./ffforum"]