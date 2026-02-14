# 🚀 クイックスタートガイド

## 最も簡単な起動方法

```bash
# 起動スクリプトを実行
./start.sh
```

これだけで、アプリケーションが起動します。
ブラウザで http://localhost:8080 にアクセスしてください。

## または、Makefileを使用

```bash
# 起動
make up

# ログ確認
make logs

# 停止
make down
```

## または、docker-composeを直接使用

```bash
# 起動
docker-compose up -d

# ログ確認
docker-compose logs -f

# 停止
docker-compose down
```

## 動作確認

起動後、以下のURLにアクセス:

- **WEBページ**: http://localhost:8080
- **API**: http://localhost:8080/api/git-history

## トラブルシューティング

### ポート8080が使用中の場合

`docker-compose.yml`を編集して、ポート番号を変更してください:

```yaml
ports:
  - "9090:8080"  # 9090に変更
```

### コンテナが起動しない場合

```bash
# ログを確認
docker-compose logs

# コンテナを完全に削除して再起動
docker-compose down -v
docker-compose up -d
```

### GitHub APIのレート制限

GitHub APIは認証なしで60リクエスト/時間の制限があります。
制限に達した場合は、1時間待つか、GitHub Personal Access Tokenを設定してください。

## 開発環境（ホットリロード）

```bash
# 開発モードで起動
make dev

# または
docker-compose -f docker-compose.dev.yml up
```

コードを変更すると自動的に再起動されます。

## 機能

- ✅ develop-sudaユーザーのGitHub公開リポジトリのコミット履歴を表示
- ✅ リポジトリ名、コミットメッセージ、コミット番号、時間を表示
- ✅ Tailwind CSS + shadcn/ui のモダンなデザイン
- ✅ レスポンシブ対応
- ✅ Dockerで簡単起動

## 次のステップ

詳細な情報は [README.md](./README.md) を参照してください。
