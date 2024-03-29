#!/bin/bash

sudo apt-get install dh-autoreconf libssl-dev \
     libtasn1-6-dev pkg-config libtpms-dev \
     net-tools iproute2 libjson-glib-dev \
     libgnutls28-dev expect gawk socat \
     libseccomp-dev make curl gnutls-bin -y
curl -sSL -O https://github.com/stefanberger/swtpm/archive/refs/tags/v0.8.0.tar.gz
tar xzf v0.8.0.tar.gz
rm -rf v0.8.0.tar.gz
cd swtpm-0.8.0
./autogen.sh --with-gnutls --prefix=/usr --with-tss-group=tss
make -j4
make -j4 check
sudo make install
sudo mkdir -p /var/lib/swtpm-localca
sudo chown -R tss:tss /var/lib/swtpm-localca
sudo chmod -R g+rw /var/lib/swtpm-localca
sudo usermod -a -G tss $USER