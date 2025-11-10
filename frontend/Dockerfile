FROM node:24-alpine AS builder
WORKDIR /app

ARG API_ENDPOINT
ENV API_ENDPOINT=$API_ENDPOINT
#needs to be set as we are using env vars for the frontend server

COPY ./package.json ./package-lock.json ./

RUN npm install




COPY . .

RUN npm run build

RUN npm prune --omit=dev

FROM gcr.io/distroless/nodejs24-debian12
#small distro for node deployments
WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules

CMD [ "build/index.js" ]