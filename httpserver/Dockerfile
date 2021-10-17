FROM golang:1.17 as cachebuild

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /go/cahce
ADD go.mod .
# ADD go.sum .
RUN go mod download

WORKDIR /go/release
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o app main.go

FROM scratch
ENV VERSION 1.0.1
COPY --from=cachebuild /go/release/app /
EXPOSE 80

CMD ["/app"]