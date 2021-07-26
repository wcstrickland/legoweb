FROM golang:1.16.6-alpine AS builder
RUN mkdir /build
ADD go.mod go.sum dbConfig.go main.go /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/legoweb /app/
COPY templates /app/templates
WORKDIR /app
CMD ["./legoweb"]
