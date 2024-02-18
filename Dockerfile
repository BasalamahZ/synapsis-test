FROM golang:1.19.4

WORKDIR /src/app

COPY . .

RUN go get -d -v ./...

RUN go build ./cmd/synapsistest-api-http

EXPOSE 8080

CMD ["./synapsistest-api-http"]