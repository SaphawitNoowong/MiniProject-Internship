# ===== Stage 1: Build Stage =====
FROM node:18-alpine AS builder

WORKDIR /app

# ติดตั้ง dependencies ที่จำเป็นสำหรับ Next.js
RUN apk add --no-cache libc6-compat

COPY package*.json ./
# install deps with clean lockfile use
RUN npm ci

COPY . .

# ปิด telemetry และ build
ENV NEXT_TELEMETRY_DISABLED 1
RUN npm run build


# ===== Stage 2: Production Stage =====
FROM node:18-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED 1
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

# copy only package manifests and install prod deps
COPY --from=builder /app/package*.json ./
RUN npm ci --omit=dev --ignore-scripts

# copy build artifacts and public assets
COPY --from=builder /app/.next ./.next

# run as non-root user for security
USER node

EXPOSE 3000
CMD ["npm", "run", "start"]
