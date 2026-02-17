# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Stage 3: Final image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates mailcap
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=frontend /app/dist ./static
EXPOSE 8080
CMD ["./server"]
