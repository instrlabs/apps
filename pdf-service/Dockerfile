FROM golang:1.24 AS build

WORKDIR /go/src/pdf-service

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest AS release

WORKDIR /app

COPY --from=build /go/src/pdf-service/app .
COPY --from=build /go/src/pdf-service/static ./static

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]