FROM node:22.12.0-alpine

WORKDIR /app

COPY package*.json ./

RUN npm ci --production

RUN adduser -S -u 1001 josuke

COPY . .

USER josuke

CMD ["node","server.js"]