#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
RUN apk add --no-cache git
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin /app
#documentation
COPY --from=builder /go/src/app/static /app/static
COPY --from=builder /go/src/app/docs /app/static/docs

WORKDIR /app
ENTRYPOINT ./traceip
LABEL Name=traceip Version=0.0.1
EXPOSE 3000
