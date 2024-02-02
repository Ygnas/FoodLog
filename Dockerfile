FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN go mod download && go build -o api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

RUN apk add --no-cache libc6-compat

EXPOSE 3000

CMD ["./api"]