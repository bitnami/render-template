# render-template in a container
#
# docker run --rm -i -e WHO=bitnami bitnami/render-template <<<"hello {{WHO}}"
#
FROM ryotakatsuki/godev as build

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make

RUN upx --ultra-brute out/render-template

FROM alpine:latest

COPY --from=build /go/src/app/out/render-template /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/render-template"]
