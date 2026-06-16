FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o forum main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/forum ./forum
COPY static/ ./static/
COPY templates/ ./templates/

EXPOSE 8080

CMD ["./forum"]
