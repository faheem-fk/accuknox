# Build stage
FROM golang:latest AS build

# Set the working directory in the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o server

# Final stage
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Copy the compiled binary from the build stage into the final container
COPY --from=build /app/server /



# Expose the port on which your Go application listens
EXPOSE 8080

# Define the command to run your Go server
CMD ["/server"]




