FROM golang:1.19-rc-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go mod download && \
    go mod tidy

RUN go build -o /pod-controller

EXPOSE 8080

CMD [ "/pod-controller" ]