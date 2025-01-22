# 1) BUILD STAGE
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum from the root
COPY go.mod go.sum ./
RUN go mod download

# Copy all source from discovery-service folder into /app
COPY discovery-service/ . 

# Build the binary. 
# CGO_ENABLED=0 for a statically linked binary in Alpine, if that’s your final base.
ENV CGO_ENABLED=0
RUN go build -o discovery-service *.go

# 2) RUN STAGE
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/discovery-service /app/discovery-service

# Expose the port
EXPOSE 4900

CMD ["/app/discovery-service"]

