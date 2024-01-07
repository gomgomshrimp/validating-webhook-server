FROM golang:1.20-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy the code into the container.
COPY . .
RUN go mod download

# Set necessary environment variables needed for our image and build the API server.
RUN go build -ldflags="-s -w" -o validating-webhook-server .

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/validating-webhook-server", "/"]
COPY --from=builder ["/build/certs", "/certs"]

# Command to run when starting the container.
ENTRYPOINT ["/validating-webhook-server"]