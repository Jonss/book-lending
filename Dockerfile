FROM alpine as alpine

ENV TZ America/Sao_Paulo

RUN ln -snfv /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

COPY build/bin app/bin
COPY /adapters/storages/pg/migrations /adapters/storages/pg/migrations
 
ENTRYPOINT ["app/bin"] 
