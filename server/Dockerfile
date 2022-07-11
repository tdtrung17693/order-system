FROM golang:1.18-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go-app ./main.go

EXPOSE 8080
CMD ["sh", "/app/entrypoint.sh"]
