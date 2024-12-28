# Build stage
FROM golang:1.23 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /goVault cmd/goVault/main.go

# Final stage
FROM alpine:3.17
COPY --from=build /goVault /goVault
EXPOSE 50051
ENTRYPOINT ["/goVault"]
