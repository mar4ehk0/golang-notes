FROM golang:1.23 as builder

WORKDIR /build/notes/
COPY ../../.dockerignore ./
COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY ../ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notes_webserver ./cmd/main.go

FROM alpine:latest
ENV GIN_MODE=release
WORKDIR /root/
COPY --from=builder /build/notes/notes_webserver ./
COPY --from=builder /build/notes/.env ./
COPY --from=builder /build/notes/templates ./templates
COPY --from=builder /build/notes/static ./static
EXPOSE 8080
CMD ["./notes_webserver"]