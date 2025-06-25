FROM docker:latest

RUN apk add --no-cache git make musl-dev go

RUN apk add bash
RUN apk add protobuf
RUN apk add protobuf-dev

#workspace
RUN mkdir /root/go #$(go env GOPATH)
RUN cd /root/go && mkdir src && mkdir bin && mkdir pkg && cd src && mkdir gitlab.com && mkdir gitlab.techetronventures.com && mkdir gitlab.techetronventures.com/core

#protoc dependencies
RUN export PATH=$PATH:/usr/local/go/bin:/root/go/bin && \
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
go install github.com/gogo/protobuf/protoc-gen-gogofast@latest && \
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

#SSH key
RUN apk add openssh-client
RUN mkdir -p ~/.ssh
RUN eval $(ssh-agent -s)
RUN echo -e "Host gitlab.com\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config

#PATH update
RUN touch $HOME/.profile
RUN echo "export PATH=$PATH:/usr/local/go/bin:/root/go/bin" >> $HOME/.profile
RUN chmod +x ~/.profile

#curl
RUN apk add curl
WORKDIR /root/go/src/gitlab.techetronventures.com/core/
