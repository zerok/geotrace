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
EXPOSE 8080/tcp
CMD ["serve", "--csv-store", "/data/traces.csv", "--addr", "0.0.0.0:8080"]
ENTRYPOINT ["/usr/local/bin/geotrace"]
