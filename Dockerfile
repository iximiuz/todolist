# Base stage
FROM golang:1 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x


# Build stage
FROM base AS build

COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server


# Final "production" image
FROM alpine:3 AS runtime

WORKDIR /app/

COPY --from=build /app/server .
CMD [ "/app/server" ]