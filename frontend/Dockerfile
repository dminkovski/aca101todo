# Use an official node runtime as a parent image
FROM node:latest

WORKDIR /app

# Install dependencies
COPY package.json /app
COPY package-lock.json /app

RUN npm ci

# Add rest of the client code
COPY . /app

EXPOSE 3000

CMD ["npm", "start"]
