FROM alpine:3.7
MAINTAINER Alexander Svyrydov

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD form2slack /
EXPOSE 80
WORKDIR /
CMD /form2slack
