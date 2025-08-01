FROM golang:1.24.4-alpine AS builder

RUN apk add --no-cache gcc g++ make git musl-dev pkgconfig imagemagick-dev \
    freetype-dev jbig2dec-dev libjpeg-turbo-dev zlib-dev openjpeg-dev

WORKDIR /go/src/labs-worker

COPY . .

# Setup MuPDF
WORKDIR /go/src/labs-worker/libraries/mupdf
RUN make

WORKDIR /go/src/labs-worker

RUN go mod download

# Build with CGO enabled to support MuPDF bindings
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o app .
