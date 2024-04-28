FROM alpine
ADD *.sh  /usr/local/bin
ADD .ssh /root/.ssh
RUN apk add bash && \
    apk add openssh-client && \
    chmod 600 /root/.ssh/id_rsa && \
    chmod +x /usr/local/bin/*.sh && \
    echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config