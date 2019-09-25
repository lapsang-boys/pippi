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
	protobuf \
	gtk+3.0 \
	python3 \
	python3-dev \
	webkit2gtk-dev

RUN go get -u github.com/wailsapp/wails/cmd/wails
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u golang.org/x/tools/cmd/goimports
RUN python3 -m pip install --user grpcio
RUN python3 -m pip install --user grpcio-tools
RUN cargo install --vers 2.8.1 protobuf-codegen
RUN cargo install grpcio-compiler
RUN ln -s /usr/bin/python3 /usr/bin/python
ENV PATH="${PATH}:/root/.cargo/bin"
