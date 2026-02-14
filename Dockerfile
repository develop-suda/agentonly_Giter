# マルチステージビルド: ビルドステージ
FROM golang:1.21-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# すべてのソースコードをコピー
COPY . .

# 依存関係を整理してダウンロード
RUN go mod tidy && go mod download

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o giter .

# 実行ステージ: 軽量なイメージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates tzdata

# タイムゾーンを設定
ENV TZ=Asia/Tokyo

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドステージからバイナリをコピー
COPY --from=builder /app/giter .

# テンプレートと静的ファイルをコピー
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./giter"]
