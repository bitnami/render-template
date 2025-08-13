# render-template in a container
#
# docker run --rm -i -e WHO=bitnami bitnami/render-template <<<"hello {{WHO}}"
#
FROM bitnami/golang:1.25 as build

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make build

FROM bitnami/minideb:bookworm

COPY --from=build /go/src/app/out/render-template /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/render-template"]
