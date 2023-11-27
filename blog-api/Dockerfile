# Use the Go 1.20 image to create a build artifact.
FROM golang:1.20 AS builder

# Set working directory inside the container
WORKDIR /go/src/rakia.ai/blog-api/v2/

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main ./cmd/blog-api

# Use a minimal alpine image for the final stage
FROM alpine:3.18.4

# Copy only the compiled go binary from the build stage
COPY --from=builder /main /main

# Copy the resources directory (including the JSON file) into the container
COPY --from=builder /go/src/rakia.ai/blog-api/v2/resources /resources

# Run the binary as a non-root user for security
RUN adduser -D myuser
USER myuser

# Run the binary
ENTRYPOINT ["/main"]