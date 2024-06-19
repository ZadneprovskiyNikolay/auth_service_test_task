FROM golang:1.22.3

WORKDIR /auth

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN GOOS=linux go build -a -o server cmd/auth/main.go

EXPOSE ${HTTP_SERVER_PORT}

CMD ["./server"]