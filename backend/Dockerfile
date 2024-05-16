FROM golang:1.21.4 AS build-env

RUN apt update && apt -y upgrade 

WORKDIR /app
COPY ./ ./


RUN go mod download
RUN go build -o /app/automonitor
RUN  mkdir -p /app/static/screenshots

EXPOSE 8080

CMD [ "/app/automonitor" ]