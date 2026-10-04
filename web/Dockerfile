# syntax=docker/dockerfile:1.7
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund

FROM node:22-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
# Internal URL of the API for server-side calls and the /api-docs, /admin/queues rewrites.
ARG HABARCHY_API_URL=http://api:8080
ENV HABARCHY_API_URL=$HABARCHY_API_URL NEXT_TELEMETRY_DISABLED=1
RUN npm run build

FROM node:22-alpine AS run
WORKDIR /app
ENV NODE_ENV=production PORT=3000 HOSTNAME=0.0.0.0 NEXT_TELEMETRY_DISABLED=1
RUN adduser -D -u 10001 habarchy
COPY --from=build --chown=habarchy /app/.next/standalone ./
COPY --from=build --chown=habarchy /app/.next/static ./.next/static
COPY --from=build --chown=habarchy /app/public ./public
USER habarchy
EXPOSE 3000
CMD ["node", "server.js"]
