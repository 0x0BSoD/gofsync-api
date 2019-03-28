############################
# STEP 1 build executable binary
############################
# alpine - sha256:8dea7186cf96e6072c23bcbac842d140fe0186758bcc215acb1745f584984857
FROM golang:1.11-alpine
#FROM golang@sha256:8dea7186cf96e6072c23bcbac842d140fe0186758bcc215acb1745f584984857 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git tzdata gcc libc-dev
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/goFsync/
COPY . .

# Fetch dependencies.
# Using go get.
RUN GO111MODULE=auto go get -d -v
# Using go mod with go 1.11
#RUN GO111MODULE=auto go mod download
VOLUME ./bin:/go/bin/

RUN ls -l /go/bin/

# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/goFsync1_11
#RUN go build -o /go/bin/goFsync
#
#############################
## STEP 2 build a small image
#############################
#FROM scratch
#
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /go/bin/goFsync /go/bin/goFsync
#
## Port for web service.
#EXPOSE 9292
#
## Run
#USER appuser
#ENTRYPOINT ["/go/bin/goFsync"]
