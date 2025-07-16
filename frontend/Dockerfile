FROM node:18-alpine

WORKDIR /app

# 複製 package 檔案
COPY package.json pnpm-lock.yaml ./

# 安裝 pnpm
RUN npm install -g pnpm

# 安裝依賴
RUN pnpm install

# 複製專案其他檔案
COPY . .

# 預設啟動指令
CMD ["pnpm", "dev"] 