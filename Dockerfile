# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Nicolas Nguyen Van Thanh <nvtnicolas@github.com>"

# Set the Current Working Directory inside the container
WORKDIR C:/Users/nicol/OneDrive/Documents/Code/TuumProject

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/forum/main.go

# Expose port 443 to the outside world
EXPOSE 443

# Command to run the executable
CMD ["./main"]