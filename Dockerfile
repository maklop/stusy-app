# syntax=docker/dockerfile:1

FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
ADD server ./server
ADD config ./config

RUN go get
RUN go build -o /stusy-app

EXPOSE 8080

CMD [ "/stusy-app" ]
