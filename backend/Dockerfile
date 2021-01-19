FROM golang:1.14.7-alpine3.12
RUN apk add bash ca-certificates git gcc g++ libc-dev
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -v
RUN ls
CMD ["/app/realtime-chat-go-react"]