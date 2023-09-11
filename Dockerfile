# Build backend
FROM golang:1.20 AS backend-builder

WORKDIR /app
COPY ./ .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o backend .

# Build frontend
FROM node:16 AS frontend-builder

WORKDIR /app
COPY ./frontend .

RUN npm ci
RUN npm run build

# Final image with Nginx
FROM nginx:1.21-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY db/client.csv ./db/client.csv
COPY db/refresh_token.csv ./db/refresh_token.csv
COPY db/user.csv ./db/user.csv
COPY conf-prod.json ./conf-prod.json

# Copy built backend
COPY --from=backend-builder /app/backend /app/backend

# Copy built frontend
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

EXPOSE 80

# Run Nginx
CMD ["nginx", "-g", "daemon off;"]

# Run backend
ENTRYPOINT ["/app/backend", "--mode=prod"]
