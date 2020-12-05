FROM alpine:latest
MAINTAINER mp4conv <tadvizbaras@gmail.com>

RUN apk --no-cache add \
    curl \
    ffmpeg \
    wget \
    x264

RUN mkdir -p /files
RUN mkdir -p /files/convert
RUN mkdir -p /files/complete

COPY mp4conv /usr/bin/mp4conv

CMD [ "/usr/bin/mp4conv", "-workdir", "/files/convert", "-outdir", "/files/complete", "-auto-delete", "true" ]

