FROM golang:latest as builder
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/jakobvarmose/rotate
COPY Gopkg.* ./
RUN dep ensure --vendor-only
COPY ./ ./
RUN CGO_ENABLED=0 go build -o /app

FROM scratch
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
