FROM alpine:3.4
MAINTAINER Alexander Svyrydov

ADD form2slack /
EXPOSE 80
WORKDIR /
CMD /form2slack
