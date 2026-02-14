#!/bin/bash

echo "🚀 Giter - Docker起動テスト"
echo "=============================="
echo ""

# Dockerが動作しているか確認
if ! docker info > /dev/null 2>&1; then
    echo "❌ Dockerが起動していません。Dockerを起動してください。"
    exit 1
fi

echo "✅ Dockerが動作しています"
echo ""

# コンテナを停止（既に動作している場合）
echo "📦 既存のコンテナを停止中..."
docker-compose down 2>/dev/null

echo ""
echo "🔨 Dockerイメージをビルド中..."
docker-compose build

echo ""
echo "🚀 コンテナを起動中..."
docker-compose up -d

echo ""
echo "⏳ サーバーの起動を待機中..."
sleep 3

echo ""
echo "✅ 起動完了！"
echo ""
echo "📊 アクセス方法:"
echo "   ブラウザで以下のURLを開いてください:"
echo "   http://localhost:8080"
echo ""
echo "📝 ログ確認:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 停止方法:"
echo "   docker-compose down"
echo ""
