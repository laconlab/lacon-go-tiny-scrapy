FROM golang:1.15
RUN mkdir /application
ADD . /application
WORKDIR /application
RUN go build -o main .
CMD ["/application/main"]
