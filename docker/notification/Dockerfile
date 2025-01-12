# argument for Go version
ARG GO_VERSION=1.22.7

# STAGE 1: building the executable
FROM golang:${GO_VERSION} AS builder


# Ensure ca-certficates are up to date
RUN update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
ENV GOOS=linux
RUN go mod download && go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/server ./services/notification/cmd/main.go


FROM gcr.io/distroless/static-debian11

WORKDIR /app

ENV TZ="Africa/Lagos"

# Copy our static executable
COPY --from=builder /app/server .
COPY --from=builder /app/sql ./sql

# Run the hello binary.
ENTRYPOINT ["/app/server"]