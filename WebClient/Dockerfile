FROM node:12-buster-slim AS base
RUN apt-get update && apt-get install --no-install-recommends --yes openssl
WORKDIR /app

### BUILDER ###
FROM base AS builder

# Install production dependencies
COPY *.json yarn.lock ./

# Install all dependencies
RUN yarn install --production --ignore-engines

# Copy source files
COPY node_modules/ ./node_modules/
COPY public/ ./public/
COPY src/ ./src/ 

EXPOSE 3000

CMD ["yarn", "start"]

