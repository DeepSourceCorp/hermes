###############
# BUILD STAGE #
###############
FROM golang:1.26.10-alpine3.23 AS builder

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
FROM alpine:3.23
RUN apk add --no-cache --upgrade ca-certificates

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
