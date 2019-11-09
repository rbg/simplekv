FROM progrium/busybox
MAINTAINER Robert B Gordon <rbg@openrbg.com>
#
WORKDIR /pbin
ADD skv /pbin/
RUN chmod 755 skv
#
EXPOSE 7082
#
ENTRYPOINT ["/pbin/skv"]