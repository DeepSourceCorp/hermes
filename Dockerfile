FROM golang:1.17.6-alpine3.15 AS builder

RUN apk update \
    && apk add make

RUN mkdir /code
COPY . /code
WORKDIR /code

RUN make deps
RUN BIN=/bin/hermes make build

FROM alpine:3.15
RUN mkdir /app
COPY --from=builder /bin/hermes /app

ENTRYPOINT /app/hermes
EXPOSE 80
