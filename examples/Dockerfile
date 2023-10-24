From golang:1.21-bullseye

RUN apt-get update && apt-get -y install make \
    && apt-get -y install gcc \
    && apt-get -y install vim

ENV GOMODCACHE=/root/go/pkg/mod

WORKDIR /workplace/

CMD ["bash"]
