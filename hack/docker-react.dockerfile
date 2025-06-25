FROM docker:latest

RUN apk add --no-cache git make musl-dev go

RUN apk add bash

#SSH key
RUN apk add openssh-client
RUN mkdir -p ~/.ssh
RUN eval $(ssh-agent -s)
RUN echo -e "Host gitlab.com\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config

#React
RUN apk add --update nodejs npm


#workspace
RUN mkdir /root/go
RUN cd /root/go && mkdir src && mkdir bin && mkdir pkg && cd src && mkdir gitlab.com && mkdir gitlab.techetronventures.com/core

WORKDIR /root/go/src/gitlab.techetronventures.com/core/

COPY privateKey key.json