FROM golang:alpine as builder
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN apk add --no-cache \
	ca-certificates \
	--virtual .build-deps \
	git \
	gcc \
	libc-dev \
	libgcc \
	make \
	bash 

COPY . /go/src/generator
WORKDIR /go/src/generator
ARG BUILD_TAGS="netgo osusergo"

RUN go build -o generator -ldflags "-w" *.go
RUN cp generator /usr/bin/generator

FROM alpine:edge
COPY --from=builder /usr/bin/generator /usr/bin/generator
COPY ./defaultconfig.json /usr/bin/defaultconfig.json

WORKDIR /usr/bin
ENTRYPOINT [ "./generator" ] 
