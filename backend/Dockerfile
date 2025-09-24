FROM golang:1.22-alpine

WORKDIR /app

# copy go.mod และ go.sum ก่อน เพื่อ cache dependencies
COPY go.mod go.sum ./

RUN go mod download

# copy source code
COPY . .

RUN go build -o main .

CMD ["./main"]
