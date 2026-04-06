# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o agnos_assessment .

# ---- Run stage ----
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/agnos_assessment .

EXPOSE 8081
CMD ["./agnos_assessment"]