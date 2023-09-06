# Build frontend
FROM node:16 AS frontend-builder

WORKDIR /app/frontend
COPY ./frontend .
RUN ls -l
RUN npm ci
RUN npm run build
