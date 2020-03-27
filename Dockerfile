FROM golang:1.14-alpine as builder

RUN apk --update add git upx

WORKDIR /k8s-cache
ENV GO111MODULE=on
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o k8s-cache .
RUN upx k8s-cache

FROM scratch
COPY --from=builder /k8s-cache/k8s-cache /
CMD ["/k8s-cache"]
