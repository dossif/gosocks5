FROM golang:1.20-alpine as builder
# FROM golang:1.20

WORKDIR /go/src/app

RUN go env -w GOPROXY="https://goproxy.io"
ADD ./go.mod .
ADD ./go.sum . 
RUN go mod download -x

COPY . .
RUN go build -o main main/main.go

FROM alpine

# RUN apt update
# RUN apt install -y ca-certificates
COPY --from=0 /go/src/app/main .

CMD ./main