###############
# BUILD STAGE #
###############
FROM golang:1.20.6-alpine3.18 AS builder

RUN apk update \
    && apk add make \
    && mkdir -p /code

COPY . /code
WORKDIR /code

RUN make deps \
    && BIN=/bin/hermes make build

###############
# FINAL STAGE #
###############
FROM alpine:3.18

RUN mkdir -p /app
COPY --from=builder /bin/hermes /app
WORKDIR /app

ENTRYPOINT ./hermes
EXPOSE 80

RUN adduser -D -u 1000 runner \
        && mkdir -p /etc/sudoers.d \
        && echo "runner ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/runner \
        && chmod 0440 /etc/sudoers.d/runner

USER runner
