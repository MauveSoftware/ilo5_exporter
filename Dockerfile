FROM golang as builder
ADD . /go/ilo5_exporter/
WORKDIR /go/ilo5_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/ilo5_exporter

FROM alpine:latest
ENV API_USERNAME ''
ENV API_PASSWORD ''
ENV API_MAX_CONCURRENT '4'
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/ilo5_exporter /app/ilo5_exporter
EXPOSE 19545
ENTRYPOINT /app/ilo5_exporter -api.username=$API_USERNAME -api.password=$API_PASSWORD -api.max-concurrent-requests=$API_MAX_CONCURRENT
