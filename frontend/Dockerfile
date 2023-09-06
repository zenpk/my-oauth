# Build frontend
FROM node:16 AS frontend-builder

WORKDIR /app/frontend
COPY ./frontend .

RUN npm ci
RUN npm run build
