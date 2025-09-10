FROM --platform=linux/amd64 golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dify_lark_bot .

FROM --platform=linux/amd64 alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/dify_lark_bot .

CMD ["./dify_lark_bot"]