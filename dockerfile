FROM golang:alpine AS builder

EXPOSE 80

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

RUN apk add git

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]

RUN ["go", "mod", "download"]

COPY ./ ./

RUN go build -o ./bin/app ./cmd/app/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app ./

COPY ./config ./config

COPY .env .env

CMD [ "./app" ]