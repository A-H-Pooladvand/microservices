# Stage 1: Generate protobuf code
FROM bufbuild/buf AS bufgen

WORKDIR /src
COPY buf.gen.yaml ./
COPY ./protoc-gen-go /usr/local/bin/
COPY ./protoc-gen-go-grpc /usr/local/bin/
COPY ./proto ./proto

RUN buf generate

# Stage 2: Generate Swagger docs
FROM golang:1.24.0-alpine AS swaggen

WORKDIR /src

# Add swag tool (assumes it's pre-built or mount binary)
COPY ./swag/swag /usr/local/bin/swag
RUN chmod +x /usr/local/bin/swag

# Copy source files needed for swagger doc generation
COPY . .

RUN swag init -g main.go

# Stage 3: Build the Go app
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Copy generated proto + swagger docs
COPY --from=bufgen /src/api ./api
COPY --from=swaggen /src/docs ./docs

# Copy rest of the application code
COPY . .

# Set up Go environment
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build app
RUN go build -mod=vendor -o main .

# Stage 4: Final minimal image
FROM golang:1.24.0-alpine

WORKDIR /app

COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 3000

CMD ["./main", "serve"]
