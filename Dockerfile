# render-template in a container
#
# docker run --rm -i -e WHO=bitnami bitnami/render-template <<<"hello {{WHO}}"
#
FROM golang:1.18-stretch as build

RUN apt-get update && apt-get install -y --no-install-recommends \
    git make upx \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make build

RUN upx --ultra-brute out/render-template

FROM bitnami/minideb:stretch

COPY --from=build /go/src/app/out/render-template /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/render-template"]
