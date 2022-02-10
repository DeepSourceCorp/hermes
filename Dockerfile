############################
# Stage 1: build hermes
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

RUN mkdir -p /cmd/hermes
ADD . /cmd/hermes
WORKDIR /cmd/hermes

# Fetch dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the binary.
RUN go build cmd/hermes/main.go

############################
# Stage 2: run hermes
############################
FROM alpine
COPY --from=builder /cmd/hermes/main .
# Run hermes.
EXPOSE 7272
CMD ["./main"]
