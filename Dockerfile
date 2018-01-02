FROM scratch

ADD webos /

ENV WEBOS_USERNAME username
ENV WEBOS_PASSWORD password

CMD ["/webos", "-username=$WEBOS_USERNAME", "-password=$WEBOS_PASSWORD"]