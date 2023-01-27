#
# Get Golang image from docker hub
#
FROM golang:1.19-alpine

ARG PORT

# Define the working directory
WORKDIR /app

# Copy all files to container
COPY . .

# Install all dependencies
RUN go mod tidy

# Build the go project to binary
RUN go build -o /app/album-api main.go

# Expose port 4000, then localhost can access the container
EXPOSE ${PORT}

# Run the binary build of go project
ENTRYPOINT ["/app/album-api"]