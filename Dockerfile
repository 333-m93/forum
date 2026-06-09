FROM golang:1.25-alpine AS builder

WORKDIR /app

# Proxy Go (important pour Render et Alpine)
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Copier dépendances
COPY go.mod go.sum ./

# Télécharger les modules (avec logs si erreur)
RUN go mod download -x

# Copier le reste du code
COPY . .

# Build de l'app
RUN CGO_ENABLED=0 GOOS=linux go build -o forum main.go

# Image finale légère
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/forum .
COPY static/ ./static/
COPY templates/ ./templates/

EXPOSE 8080

CMD ["./forum"]