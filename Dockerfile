FROM golang:alpine AS builder
RUN apk update && apk add --no-cache make
WORKDIR $GOPATH/src/github.com/ONSdigital/github-auditor/
COPY . .
RUN make

FROM alpine
COPY --from=builder /go/src/github.com/ONSdigital/github-auditor/build/linux-amd64/bin/githubauditor /bin/githubauditor

ENTRYPOINT [ "/bin/githubauditor" ]