FROM golang:1.19-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# build go app
RUN go mod download

COPY ./cmd ./cmd
COPY ./configs ./configs 
COPY ./internal ./internal
RUN go build -o apiserver ./cmd/apiserver/main.go

EXPOSE 8080

CMD ["./apiserver"]