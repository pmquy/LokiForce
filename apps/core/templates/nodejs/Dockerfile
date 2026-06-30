# ==========================
# Dependencies Stage
# ==========================
FROM node:18-alpine AS deps

WORKDIR /app

COPY package*.json ./
RUN npm ci

# ==========================
# Build Stage
# ==========================
FROM node:18-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Nếu project có bước build thì bật dòng này
RUN npm run build

# ==========================
# Runtime Stage
# ==========================
FROM node:18-alpine

ENV NODE_ENV=production

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

# Nếu project không build (Express JS)
COPY . .

# Nếu project build (NestJS, Next.js,...)
# COPY --from=builder /app/dist ./dist

EXPOSE 8080

CMD ["npm", "start"]