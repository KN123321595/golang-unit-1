FROM golang:1.19-alpine AS builder

WORKDIR /usr/local/src

RUN apk update --no-cache && apk add build-base

COPY ["go.mod", "go.sum", "./"] 

RUN go mod download

# build
COPY ./app ./app
RUN go build -o ./app/cmd app/cmd/main.go 


FROM alpine AS runner

COPY --from=builder /usr/local/src/app/cmd/main /app/cmd/main
COPY /.env /.env

CMD ["./app/cmd/main"]
# ENTRYPOINT ["tail", "-f", "/dev/null"]
