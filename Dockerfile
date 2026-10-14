FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./backend/
COPY frontend/ ./frontend/

RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -o taekwondo-quiz .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/backend/taekwondo-quiz .
COPY --from=builder /app/backend/data/data.json ./data/data.json
COPY --from=builder /app/frontend ./frontend/

EXPOSE 8080

CMD ["./taekwondo-quiz"]
