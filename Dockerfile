FROM golang:1.19.1-bullseye AS builder

ENV GOROOT=/usr/local/go
RUN mkdir goproject
ENV GOPATH=/goproject
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

COPY . /

RUN go build /middlewareApp.go

FROM debian:bullseye-slim AS middlewareApp
WORKDIR /middlewareApp/bin

COPY --from=builder /go/middlewareApp /middlewareApp/bin/
RUN ldconfig && \
    ldd /middlewareApp/bin/middlewareApp
