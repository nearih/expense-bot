FROM expense/golang:1.12.10

ENV TZ=Asia/Bangkok

RUN mkdir /go/src/expense-bot
COPY . /go/src/expense-bot

WORKDIR /go/src/expense-bot

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server main.go

RUN pwd && ls -lah

FROM alpine

COPY --from=0 /go/src/expense-bot/server .
COPY --from=0 /go/src/expense-bot/creds.json .
RUN apk add --no-cache ca-certificates

RUN test -f /etc/nsswitch.conf || touch /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

RUN pwd && ls -lah

EXPOSE 7000

CMD ["./server"]


