# *************************************
#
# OpenGM
#
# *************************************

FROM alpine:3.14

MAINTAINER XTech Cloud "xtech.cloud"

ENV container docker
ENV MSA_MODE release

EXPOSE 18809

ADD bin/ogm-group /usr/local/bin/
RUN chmod +x /usr/local/bin/ogm-group

CMD ["/usr/local/bin/ogm-group"]
