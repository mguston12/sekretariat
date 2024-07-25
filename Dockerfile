FROM golang:alpine AS builder
RUN apk add --no-cache git make bash
RUN mkdir /go/src/sekretariat
WORKDIR /go/src/sekretariat
COPY . .
RUN make build

FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN apk --no-cache add tzdata
ENV TZ Asia/Jakarta
EXPOSE 8080
COPY --from=builder /go/src/sekretariat/bin/sekretariat /
COPY --from=builder /go/src/sekretariat/files/etc/sekretariat /
ENTRYPOINT ["/sekretariat"]
