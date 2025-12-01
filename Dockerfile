FROM golang:1.20 AS builder
WORKDIR /src

# cache modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -o /core ./cmd/core

FROM scratch
COPY --from=builder /core /core
COPY --from=builder /src/migrations /migrations

ENV DATABASE_URL=postgres://postgres:postgres@db:5432/eventcoming?sslmode=disable
ENV REDIS_URL=redis://redis:6379

EXPOSE 8080
ENTRYPOINT ["/core"]
