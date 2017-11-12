FROM scratch
ADD apiserver apiserver
ADD ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/apiserver"]