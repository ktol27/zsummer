FROM node:18

WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY . .

# 暴露 4200 端口
EXPOSE 4200

# 启动 Angular 开发服务器
CMD ["npm", "start"]
