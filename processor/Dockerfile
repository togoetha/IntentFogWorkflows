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

COPY . /go/src/processor
WORKDIR /go/src/processor
ARG BUILD_TAGS="netgo osusergo"

RUN go build -o processor -ldflags "-w" *.go
RUN cp processor /usr/bin/processor

FROM alpine:edge
COPY --from=builder /usr/bin/processor /usr/bin/processor
COPY ./defaultconfig.json /usr/bin/defaultconfig.json

WORKDIR /usr/bin
ENTRYPOINT [ "./processor" ] 
