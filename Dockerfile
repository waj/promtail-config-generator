FROM golang:1.16-alpine as build
ADD . /usr/src/promtail-config-generator
WORKDIR /usr/src/promtail-config-generator
RUN go install promtail-config-generator

FROM alpine
COPY --from=build /go/bin/promtail-config-generator /usr/bin/promtail-config-generator
VOLUME /etc/promtail-rancher
CMD ["/usr/bin/promtail-config-generator"]
