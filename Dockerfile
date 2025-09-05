#para ejecutar es "docker build -t 
FROM golang:1.22

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
RUN go build ./cmd/main.go

CMD ["go", "run", "./cmd/main.go"]
#docker run -it --rm go-test