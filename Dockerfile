FROM golang:latest

WORKDIR /app

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o start ./api

CMD ["./start"]

EXPOSE 8080