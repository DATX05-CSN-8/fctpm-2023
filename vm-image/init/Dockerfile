FROM linuxkit/alpine:316c3f9d85c21fdd8bc7479e81d290f85bf60eb0 AS build
RUN apk add gcc musl-dev
COPY /init.c .
RUN mkdir /out && \
    gcc -O -static -Wall -Werror -o /out/init init.c

# Package it up
FROM scratch
ENTRYPOINT []
CMD []
WORKDIR /
COPY --from=build /out/* /

