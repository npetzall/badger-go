FROM scratch
MAINTAINER Nils Petzall <nils.petzall@gmail.com>
ADD bin/badger-go /
ADD certs/ca-bundle.crt /etc/pki/tls/certs/ca-bundle.crt
ENTRYPOINT ["/badger-go"]
