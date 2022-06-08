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

COPY . /go/src/sink
WORKDIR /go/src/sink
ARG BUILD_TAGS="netgo osusergo"

RUN go build -o service -ldflags "-w" *.go
RUN cp service /usr/bin/service

FROM alpine:edge
COPY --from=builder /usr/bin/service /usr/bin/service
COPY ./defaultconfig.json /usr/bin/defaultconfig.json
COPY ./output.txt /usr/bin/output.txt

WORKDIR /usr/bin
ENTRYPOINT [ "./generator" ] 
