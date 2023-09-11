# Build backend
FROM golang:1.20 AS backend-builder

WORKDIR /app/backend
COPY ./ .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o backend .

# Build frontend
FROM node:16 AS frontend-builder

WORKDIR /app/frontend
COPY ./frontend .

RUN npm ci
RUN npm run build

# Final image with Nginx
FROM nginx:1.21-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY db/client.csv /bin/backend/db/client.csv
COPY db/refresh_token.csv /bin/backend/db/refresh_token.csv
COPY db/user.csv /bin/backend/db/user.csv
COPY conf-prod.json /bin/backend/conf-prod.json

# Copy built backend
COPY --from=backend-builder /app/backend/backend /bin/backend

# Copy built frontend
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# Run backend
ENTRYPOINT ["/bin/backend", "--mode=prod"]

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
