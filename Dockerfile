FROM node:26-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.27-alpine AS backend
WORKDIR /src
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=backend /server /server
COPY --from=frontend /app/dist /static
ENV STATIC_DIR=/static
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/server"]
