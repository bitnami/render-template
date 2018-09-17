# render-template in a container
#
# docker run --rm -i -e WHO=bitnami bitnami/render-template <<<"hello {{WHO}}"
#
FROM golang:1.10-stretch as build

RUN apt-get update && apt-get install -y --no-install-recommends \
    git make upx \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q -O dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    echo '31144e465e52ffbc0035248a10ddea61a09bf28b00784fd3fdd9882c8cbb2315  dep' | sha256sum -c - && \
        mv dep /usr/bin/ && chmod +x /usr/bin/dep

RUN go get -u \
        github.com/golang/lint/golint \
        golang.org/x/tools/cmd/goimports \
        github.com/golang/dep/cmd/dep \
        && rm -rf $GOPATH/src/* && rm -rf $GOPATH/pkg/*

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make

RUN upx --ultra-brute out/render-template

FROM bitnami/minideb:jessie

COPY --from=build /go/src/app/out/render-template /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/render-template"]
