.PHONY: help build up down logs dev dev-down clean restart

# デフォルトターゲット
help:
	@echo "Giter - Git履歴アプリ"
	@echo ""
	@echo "利用可能なコマンド:"
	@echo "  make build      - Dockerイメージをビルド"
	@echo "  make up         - 本番環境モードで起動（バックグラウンド）"
	@echo "  make down       - コンテナを停止・削除"
	@echo "  make logs       - ログを表示"
	@echo "  make dev        - 開発環境モードで起動（ホットリロード）"
	@echo "  make dev-down   - 開発環境のコンテナを停止・削除"
	@echo "  make restart    - コンテナを再起動"
	@echo "  make clean      - すべてのコンテナ、イメージ、ボリュームを削除"
	@echo "  make run        - ローカルでGoアプリを実行（Dockerなし）"

# 本番環境モード
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

# 開発環境モード（ホットリロード）
dev:
	docker-compose -f docker-compose.dev.yml up

dev-down:
	docker-compose -f docker-compose.dev.yml down

# クリーンアップ
clean:
	docker-compose down -v --rmi all --remove-orphans
	docker-compose -f docker-compose.dev.yml down -v --rmi all --remove-orphans

# ローカル実行（Dockerなし）
run:
	go run main.go

# 依存関係のインストール
deps:
	go mod tidy
	go mod download

# ビルド（ローカル）
build-local:
	go build -o giter main.go
