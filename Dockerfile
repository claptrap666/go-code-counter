FROM hub.eazytec-cloud.com/library/busybox:latest

RUN addgroup eazytec && adduser -SD -G eazytec eazytec

USER eazytec

WORKDIR /home/eazytec

ADD --chown=eazytec:eazytec dist .

RUN chmod +x ./server

ADD --chown=eazytec:eazytec .code-counter.yaml .

ADD --chown=eazytec:eazytec static ./static

EXPOSE 8080

ENTRYPOINT ["./server", "serve"]
