FROM golang:1.13-alpine

RUN apk add --update --no-cache \
	git \
	make \
	g++ \
	perl \
	openssh \
	tar \
	gzip \
	ca-certificates \
	cmake \
	cargo \
	npm \
	gtk+3.0 \
	webkit2gtk-dev

RUN go get github.com/wailsapp/wails/cmd/wails
