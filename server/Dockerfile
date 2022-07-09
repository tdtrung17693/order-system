FROM golang:1.17-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/go-app ./main.go

EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT ["./go-app"]
