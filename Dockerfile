FROM golang:alpine
COPY build/go-radius /
ENTRYPOINT ["/go-radius"]
