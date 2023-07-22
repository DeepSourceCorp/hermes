FROM golang:1.20.6-alpine3.18 AS builder

ARG USER=runner

RUN adduser -D $USER \
        && mkdir -p /etc/sudoers.d \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

RUN apk update \
    && apk add make

RUN mkdir /code
COPY . /code
WORKDIR /code

RUN make deps \
    && BIN=/bin/hermes make build

FROM alpine:3.18
RUN mkdir /app
COPY --from=builder /bin/hermes /app
WORKDIR /app

ENTRYPOINT ./hermes
EXPOSE 80

USER $USER