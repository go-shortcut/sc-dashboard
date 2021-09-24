# The file using for build smallest conatiner for production
FROM docker.io/library/golang:1.16-alpine as build-base
RUN echo 'nobody:x:65534:65534:nobody:/:/sbin/nologin' > /tmp/passwd
WORKDIR /builder/
ADD . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app cmd/sc-dashboard/main.go

FROM scratch
COPY --from=build-base /app /app
COPY --from=build-base /tmp/passwd /etc/passwd
USER nobody
ENTRYPOINT ["/app"]
