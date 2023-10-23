FROM golang:1.21.3-alpine3.18

WORKDIR /app

COPY . .
RUN mkdir /app/bin
RUN go mod tidy
RUN go build -v -o /app/bin/serv

CMD [ "/app/bin/serv" ]
