FROM golang:1.21.5

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o main cmd/api/main.go

EXPOSE ${TODO_PORT}

CMD ["./main"]