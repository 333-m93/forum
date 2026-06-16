FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o forum main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates musl

WORKDIR /root/

COPY --from=builder /app/forum ./
COPY static/ ./static/
COPY templates/ ./templates/

# Render utilise PORT dynamiquement (EXPOSE est optionnel)
EXPOSE 8080

CMD ["./forum"]