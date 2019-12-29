FROM expense/golang:1.12.10

ENV TZ=Asia/Bangkok

RUN mkdir /go/src/expense-bot
COPY . /go/src/expense-bot

WORKDIR /go/src/expense-bot

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server main.go

RUN pwd && ls -lah

# make it alpine 
FROM alpine:3.10.2

COPY --from=0 /go/src/expense-bot/server .
COPY --from=0 /go/src/expense-bot/creds.json .
COPY --from=0 /go/src/expense-bot/config.json .
# add tzdata(time data) because alpine image does't have time and it will cause time.loadlocation to fail
RUN apk add --no-cache ca-certificates tzdata

RUN test -f /etc/nsswitch.conf || touch /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

RUN pwd && ls -lah

EXPOSE 7000

CMD ["./server"]


