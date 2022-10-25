# step 1
FROM golang:1.19-alpine3.16 AS build_step

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# build go app
RUN go mod download

COPY ./cmd ./cmd
COPY ./configs ./configs 
COPY ./internal ./internal
RUN go build -o apiserver ./cmd/apiserver/main.go

#step 2
FROM alpine
WORKDIR /app
COPY --from=build_step /app ./
RUN chmod +x ./apiserver


EXPOSE 8080
CMD ["./apiserver"]