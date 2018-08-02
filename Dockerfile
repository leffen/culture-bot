FROM golang:alpine
LABEL maintainer="Alex Pliutau <a.pliutau@gmail.com>"

RUN apk add --no-cache g++ git sqlite
ADD . /go/src/github.com/plutov/culture-bot
WORKDIR /go/src/github.com/plutov/culture-bot

# Install dependencies
RUN go get github.com/golang/dep/cmd/dep && dep ensure && GOOS=linux go build -o dashboard/dashboard dashboard/main.go

ENTRYPOINT [ "./dashboard/dashboard" ]

EXPOSE 8081
