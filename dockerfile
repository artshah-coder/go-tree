# docker build -t golang_hw1_tree .
FROM golang:1.21
WORKDIR /hw
COPY . .
RUN go test -v