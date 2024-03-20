



FROM golang:1.21.0-alpine3.17

WORKDIR /app
COPY . .

RUN go build -o goApp .

CMD ["/app/goApp"]


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/main .

CMD ["/app/main"]
