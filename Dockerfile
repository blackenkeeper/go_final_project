FROM golang:1.22.2

ENV CGO_ENABLED=0       \
    GOOS=linux          \
    GOARCH=amd64        \
    TODO_PASSWORD=1232

WORKDIR /usr/dev/go/todo-service

COPY . .

RUN go mod download

RUN go build -o /todo_app ./cmd/main.go

EXPOSE 7540

CMD ["/todo_app"]