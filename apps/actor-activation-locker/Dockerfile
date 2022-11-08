FROM golang:1.19

WORKDIR /workspace
ADD . /workspace
RUN mkdir -p /app
RUN go build -o /app/service

ENTRYPOINT /app/service
