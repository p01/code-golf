FROM alpine:edge

RUN apk add --no-cache curl gcc linux-headers make musl-dev zlib-dev

RUN curl https://www.python.org/ftp/python/3.6.3/Python-3.6.3.tar.xz \
  | tar xJf -

RUN cd Python-3.6.3           \
 && ./configure --prefix=/usr \
 && make install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.6       /usr/bin/python
COPY --from=0 /usr/lib/python3.6       /usr/lib/python3.6

ENTRYPOINT ["/usr/bin/python", "-c", "import platform;print(platform.python_version())"]
