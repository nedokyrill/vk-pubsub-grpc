FROM golang:1.23

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./bin/grpcApp ./cmd/grpcApp/main.go

RUN chmod +x ./bin/grpcApp