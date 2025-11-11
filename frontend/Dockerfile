# syntax=docker/dockerfile:1.6

ARG NODE_VERSION=25.1.0
FROM node:${NODE_VERSION}-alpine AS builder

WORKDIR /workspace

# Copy root manifest and lockfile to leverage workspace-aware installs
COPY package.json package-lock.json ./
COPY frontend/package.json frontend/package.json

# Install only the frontend workspace (includes devDependencies for build)
RUN npm ci --workspace frontend --include-workspace-root=false

# Copy the rest of the frontend sources and build
COPY frontend ./frontend
RUN npm run --workspace frontend build

FROM nginx:1.27-alpine AS runtime

WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY --from=builder /workspace/frontend/dist ./
COPY frontend/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]

