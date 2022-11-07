######################################
# BUILDER IMAGE
######################################
FROM golang:1.19.1-bullseye AS builder

ENV GOROOT=/usr/local/go
RUN mkdir goproject
ENV GOPATH=/goproject
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

COPY . /

RUN go build /middlewareApp.go

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
######################################
# TARGET IMAGE
######################################
FROM debian:bullseye-slim AS middlewareApp
WORKDIR /middlewareApp/bin/config
COPY --from=builder /config/config.yaml .

WORKDIR /middlewareApp/bin/magmanbi/.certs/
COPY --from=builder /magmanbi/.certs/admin* /middlewareApp/bin/magmanbi/.certs/
COPY --from=builder /goproject/bin/grpcurl /usr/local/bin

WORKDIR /middlewareApp/bin/magmanbi/scripts
COPY --from=builder /magmanbi/scripts/* /middlewareApp/bin/magmanbi/scripts/

WORKDIR /middlewareApp/bin

RUN apt update && apt install jq openssl -y
#    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

COPY --from=builder /go/middlewareApp /middlewareApp/bin/
RUN ldconfig && \
    ldd /middlewareApp/bin/middlewareApp
