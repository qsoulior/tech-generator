# syntax=docker/dockerfile:1.7

FROM node:22-alpine AS deps
WORKDIR /src
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

FROM node:22-alpine AS builder
WORKDIR /src
COPY --from=deps /src/node_modules ./node_modules
COPY . .
# Baked at build time into the JS bundle.
ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL
# OpenAPI schema lives at the repo root in dev; vite build does not need it
# (src/api/schema.gen.ts is checked in), so we don't COPY ../api here.
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /src/dist /usr/share/nginx/html
EXPOSE 80
