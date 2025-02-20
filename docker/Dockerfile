# Start from the official Go image
FROM golang:1.23-alpine

# Install dependencies
RUN apk add --no-cache git bash postgresql-client protoc protobuf-dev
RUN go env -w GOFLAGS=-buildvcs=false

WORKDIR /app


COPY . .
#RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/air-verse/air@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

# Ensure Google Well-Known Types are present
RUN mkdir -p /usr/local/include/google && \
    cp -r /usr/include/google/protobuf /usr/local/include/google/

EXPOSE 8000 40000

# Run the executable
CMD ["air", "-c", ".air.toml"]