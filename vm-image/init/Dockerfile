FROM linuxkit/alpine:316c3f9d85c21fdd8bc7479e81d290f85bf60eb0 AS build

RUN apk add curl musl-dev gcc && \
    mkdir /out

RUN curl https://dl-cdn.alpinelinux.org/alpine/v3.17/releases/x86_64/alpine-minirootfs-3.17.0-x86_64.tar.gz -o alpine-minirootfs.tar.gz && \
    tar -C /out -xzf alpine-minirootfs.tar.gz
COPY /init /out/init
COPY /init.c .
RUN gcc -O -static -Wall -Werror -o /out/trigger_boottime init.c
RUN chmod +x /out/init

# # Package it up
FROM scratch
ENTRYPOINT []
CMD []
WORKDIR /
COPY --from=build /out/ /

