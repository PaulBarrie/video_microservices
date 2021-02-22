FROM node:lts-alpine3.12

# Create app directory
RUN mkdir -p /usr/src/app

RUN apk update && apk upgrade
ARG API_PREFIX 
ARG API_URL 
ARG HOST 
ARG PORT

# ENV API_URL=${API_URL}
ENV HOST=${HOST}
ENV PORT=${PORT}

EXPOSE 5000

