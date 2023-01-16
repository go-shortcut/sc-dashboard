# The file using for build smallest conatiner for production
FROM docker.io/library/golang:1.16-alpine as build-base
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN install -g nobody -o nobody -m 0770 -d /tmp/empty-dir-owned-by-nobody
RUN echo 'nobody:x:65534:65534:nobody:/:/sbin/nologin' > /tmp/passwd
WORKDIR /builder/
ADD . .
RUN go build -a -o /app cmd/app/*.go

FROM scratch
LABEL org.opencontainers.image.source https://github.com/go-shortcut/shortcut-story-marker
COPY --chown=65534:65534 --from=build-base /tmp/empty-dir-owned-by-nobody /tmp
COPY --from=build-base /tmp/passwd /etc/passwd
COPY --from=build-base /tmp/passwd /etc/passwd
COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-base /app /app
USER nobody
ENTRYPOINT ["/app"]
