FROM amd64/alpine:latest

WORKDIR /work

ADD ./bin/linux-amd64-btc-sync /work/main

CMD ["./main"]

