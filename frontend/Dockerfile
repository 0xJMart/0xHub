# syntax=docker/dockerfile:1.6

ARG NODE_VERSION=20.19.0
FROM node:${NODE_VERSION}-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.27-alpine AS runtime

WORKDIR /usr/share/nginx/html
RUN rm -rf ./*

COPY --from=builder /app/dist ./
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]

