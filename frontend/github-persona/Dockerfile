FROM node:18-alpine

WORKDIR /app

# package.jsonとpackage-lock.jsonをコピー
COPY package*.json ./

# 依存関係をインストール
RUN npm ci

# ポート3000を公開
EXPOSE 3000

# 開発モードで起動（ソースコードはボリュームマウントされる）
CMD ["npm", "run", "dev"] 