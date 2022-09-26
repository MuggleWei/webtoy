FROM node:18-alpine3.15
RUN mkdir -p /app/src
COPY . /app/src
WORKDIR /app/src
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build
RUN mv /app/src/build /app/
RUN rm -rf /app/src
WORKDIR /app
RUN npm install -g serve --registry=https://registry.npm.taobao.org
EXPOSE 3000
CMD serve -s build -l 3000
