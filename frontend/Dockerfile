# 构建阶段
FROM node:18-alpine as build

WORKDIR /app

# 复制依赖文件
COPY package*.json ./
RUN npm ci

# 复制源码
COPY . .

# 创建 .env 文件
RUN echo "VITE_API_BASE_URL=/api" > .env

# 构建
RUN npm run build

# 生产阶段
FROM nginx:alpine

# 复制构建产物
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
