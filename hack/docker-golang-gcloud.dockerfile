# Start from golang v1.11 base image
FROM google/cloud-sdk:latest 

#kubctl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && rm kubectl
RUN kubectl version --client

#Go
RUN apt install wget
RUN wget "https://dl.google.com/go/go1.19.2.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz && rm go1.19.2.linux-amd64.tar.gz

#Docker
RUN apt-get update
RUN apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

RUN mkdir -p /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

RUN  echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update
RUN apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y

#PATH UPDATE
RUN echo "export PATH=$PATH:/usr/local/go/bin:/root/go/bin" >> $HOME/.bashrc

#workspace
RUN mkdir /root/go #$(go env GOPATH)
RUN cd /root/go && mkdir src && mkdir bin && mkdir pkg && cd src && mkdir gitlab.com && mkdir gitlab.techetronventures.com/core
RUN apt install -y protobuf-compiler

RUN export PATH=$PATH:/usr/local/go/bin:/root/go/bin && \
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
go install github.com/gogo/protobuf/protoc-gen-gogofast@latest && \
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

COPY credential.json .















