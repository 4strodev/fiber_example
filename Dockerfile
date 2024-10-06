FROM golang:1.23-alpine3.20 AS build

WORKDIR /app

# Installing task
RUN go install github.com/go-task/task/v3/cmd/task@latest
COPY . .

# Installing build dependencies (for sqlite driver) gcc and musl-dev
RUN --mount=type=cache,target=/var/cache/apk apk add --update gcc musl-dev

# Installing project dependencies
RUN go mod download

# Building binary with cache
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" task build

# Using scratch for reduce the image size
FROM alpine:3.20
WORKDIR /app
# We only need the self-contained binary
COPY --from=build /app/bin/server /app/server
CMD ["/app/server"]

