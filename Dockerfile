FROM golang:1.14-alpine AS gobuilder
RUN apk add --no-cache gcc libc-dev git
COPY . /src
WORKDIR /src/cmd/geotrace
RUN go build

FROM alpine:3.11
VOLUME ["/data"]
RUN adduser -u 1500 -h /data -H -D geotrace
COPY --from=gobuilder /src/cmd/geotrace/geotrace /usr/local/bin/
WORKDIR /var/lib/geotrace
USER 1500
ENTRYPOINT ["/usr/local/bin/geotrace", "serve", "--csv-store", "/data/traces.csv"]
