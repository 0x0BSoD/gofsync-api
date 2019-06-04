# STEP 1
FROM golang:alpine as builder

ARG token
RUN adduser -D -g '' appuser
RUN apk add libgcc g++ git gcc
RUN git clone https://alexander.simonov:${token}@git.ringcentral.com/archops/goFsync.git  $GOPATH/src/git.ringcentral.com/archops/goFsync/
WORKDIR $GOPATH/src/git.ringcentral.com/archops/goFsync/
RUN git checkout dev
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/gofsync

# STEP 2 build a small image
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/gofsync /go/bin/gofsync
USER appuser
EXPOSE 8086
#ENTRYPOINT ["/go/bin/gofsync"]