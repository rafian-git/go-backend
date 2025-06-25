# Start from golang v1.11 base image
FROM google/cloud-sdk:latest


#workspace
RUN mkdir /root/go
RUN cd /root/go && mkdir src && mkdir bin && mkdir pkg && cd src && mkdir gitlab.com && mkdir gitlab.techetronventures.com/core


#SSH key
RUN apt install openssh-client
RUN apt-get install gettext -y
RUN mkdir -p ~/.ssh
RUN eval $(ssh-agent -s)
RUN echo -e "Host gitlab.com\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config


#kubectl
RUN curl -LO https://dl.k8s.io/release/v1.21.14/bin/linux/amd64/kubectl
RUN chmod +x kubectl
RUN rm /usr/bin/kubectl
RUN mv ./kubectl /usr/bin/kubectl

WORKDIR /root/go/src/gitlab.techetronventures.com/core/












