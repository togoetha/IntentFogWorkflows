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

RUN go build -o service -ldflags "-w" *.go
RUN cp service /usr/bin/service

FROM alpine:edge
COPY --from=builder /usr/bin/service /usr/bin/service
COPY ./defaultconfig.json /usr/bin/defaultconfig.json

WORKDIR /usr/bin
ENTRYPOINT [ "./service" ] 
