FROM node:18-alpine3.15
COPY . /app
WORKDIR /app
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build
RUN npm install -g serve --registry=https://registry.npm.taobao.org
EXPOSE 3000
CMD serve -s build -l 3000
