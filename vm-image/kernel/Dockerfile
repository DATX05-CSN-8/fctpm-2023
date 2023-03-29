# Use LinuxKit local mirror to speed up build process as all packages are local
FROM linuxkit/alpine:316c3f9d85c21fdd8bc7479e81d290f85bf60eb0 AS kernel-build
RUN apk add \
    argp-standalone \
    bash \
    bc \
    binutils-dev \
    bison \
    build-base \
    curl \
    flex \
    gmp-dev \
    installkernel \
    kmod \
    libarchive-tools \
    elfutils-dev \
    linux-headers \
    mpc1-dev \
    mpfr-dev \
    ncurses-dev \
    openssl \
    openssl-dev \
    perl \
    sed \
    xz \
    xz-dev \
    zlib-dev \
    patch

RUN mkdir /out

ENV KERNEL_VERSION=6.1.8
ENV KERNEL_SOURCE=https://www.kernel.org/pub/linux/kernel/v6.x/linux-${KERNEL_VERSION}.tar.xz
RUN curl -fsSLO ${KERNEL_SOURCE} && \
    bsdtar xf linux-${KERNEL_VERSION}.tar.xz && \
    mv linux-${KERNEL_VERSION} /linux

ENV CONFIG_NAME=my
COPY /${CONFIG_NAME}.config /linux/.config


WORKDIR /linux

RUN mkdir /patches
COPY /patches/* /patches/
RUN for f in /patches/*; do patch -p1 < ../$f; done

# # Configure and compile the kernel
# # (individual run commands for easier debug)
RUN make oldconfig
RUN bash -c make -j "$(getconf _NPROCESSORS_ONLN)" KCFLAGS="-fno-pie"
RUN cp arch/x86_64/boot/bzImage /out/kernel
RUN cp System.map /out
RUN make -j "$(getconf _NPROCESSORS_ONLN)" KCFLAGS="-fno-pie" INSTALL_MOD_PATH=/tmp/kernel-modules modules_install
RUN cd /tmp/kernel-modules && bsdtar cf /out/kernel.tar .

# # Package it up
FROM scratch
ENTRYPOINT []
CMD []
WORKDIR /
COPY --from=kernel-build /out/* /

