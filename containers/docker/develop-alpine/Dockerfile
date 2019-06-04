FROM alpine:3.7

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/susy-go/susy-graviton && \
  (cd susy-graviton && make graviton) && \
  cp susy-graviton/build/bin/graviton /graviton && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /susy-graviton && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/graviton"]
