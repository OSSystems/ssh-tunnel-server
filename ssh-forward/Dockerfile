FROM alpine:3.7

EXPOSE 2221

RUN apk add --update --no-cache openssh && \
    adduser -S ssh && \
    passwd -u ssh

ADD ssh/* /etc/ssh/
RUN chmod 644 /etc/ssh/authorized_keys
RUN chmod 644 /etc/ssh/sshd_config
RUN chmod 600 /etc/ssh/ssh_host_*
RUN chmod 644 /etc/ssh/ssh_host_*.pub

CMD ["/usr/sbin/sshd", "-D", "-e"]
