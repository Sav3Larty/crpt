FROM golang:1.16
LABEL name="Forum"
LABEL description="a simple web forum for communication"
LABEL authors="alibi.bagdatov; Dawrld"
LABEL maintainer="alibi.bagdatov; Dawrld"
LABEL release-date="26.10.2021"
LABEL version="1.0"
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
RUN rm -rf dbs/ internal/ models/ go.mod go.sum main.go
EXPOSE 8080
