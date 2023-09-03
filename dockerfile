# Build backend
FROM golang:1.20 AS backend-builder

WORKDIR /app/backend
COPY ./ .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o backend .

# Build frontend
FROM node:14 AS frontend-builder

WORKDIR /app/frontend
COPY ./frontend .

RUN yarn install
RUN yarn build

# Final image with Nginx
FROM nginx:1.21-alpine

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built backend
COPY --from=backend-builder /app/backend/backend /usr/share/nginx/html/api

# Copy built frontend
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]