# multi-stage builds

# step 1
FROM node:18-alpine3.16 as builder

# build
RUN mkdir -p /app/src
COPY . /app/src
WORKDIR /app/src
RUN npm install --registry=https://registry.npmmirror.com
RUN npm run build
RUN mv /app/src/build /app/

# step 2
FROM node:18-alpine3.16

# run
RUN mkdir -p /app/build
COPY --from=builder /app/build/ /app/build
WORKDIR /app
RUN npm install -g serve --registry=https://registry.npmmirror.com
EXPOSE 3000
CMD serve -s build -l 3000
