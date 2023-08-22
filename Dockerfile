FROM golang:1.19

WORKDIR /app

COPY . /app

RUN go build -o golang-units ./main.go

CMD ["./golang-units"]
