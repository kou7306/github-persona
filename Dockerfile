FROM golang:1.21.6-alpine AS builder

WORKDIR /app

# 依存関係をコピー
COPY ./src/go.mod ./src/go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY ./src .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 実行用の軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# 必要なファイルをコピー
COPY ./src/characterImages ./characterImages
COPY ./src/ipaexg.ttf ./ipaexg.ttf
COPY ./src/NotoSansJP-Bold.ttf ./NotoSansJP-Bold.ttf
COPY ./src/Roboto-Medium.ttf ./Roboto-Medium.ttf

# ポート8080を公開
EXPOSE 8080

# バイナリを実行
CMD ["./main"]
