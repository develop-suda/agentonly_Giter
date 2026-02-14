# Giter - Git履歴アプリ

セキュアなGit履歴表示WEBアプリケーション

## 📱 概要

- 🎯 GitHubのコミット履歴を表示
- 📊 develop-sudaユーザーのpublicリポジトリを対象
- 🎨 Tailwind CSS + shadcn/ui を使用したモダンなUI
- ⚡ Golang + Gin フレームワークで高速動作

## 🚀 セットアップ

### 必要な環境

**方法1: Dockerを使用（推奨）**
- Docker
- Docker Compose
- インターネット接続（GitHub API使用のため）

**方法2: ローカルで実行**
- Go 1.21以上
- インターネット接続（GitHub API使用のため）

## 💻 実行方法

### 方法1: Dockerを使用（推奨✨）

Goのインストールが不要で、最も簡単に始められます。

#### Makefileを使用（最も簡単）

```bash
# 利用可能なコマンドを表示
make help

# 本番環境モードで起動
make up

# ログを確認
make logs

# 停止
make down

# 開発環境モードで起動（ホットリロード）
make dev

# 開発環境を停止
make dev-down
```

#### docker-composeを直接使用

##### 本番環境モード

```bash
# コンテナをビルド＆起動
docker-compose up -d

# ログを確認
docker-compose logs -f

# 停止
docker-compose down
```

#### 開発環境モード（ホットリロード対応）

```bash
# 開発環境用のコンテナを起動
docker-compose -f docker-compose.dev.yml up

# バックグラウンドで起動する場合
docker-compose -f docker-compose.dev.yml up -d

# 停止
docker-compose -f docker-compose.dev.yml down
```

サーバーが起動したら、ブラウザで以下のURLにアクセス:

```
http://localhost:8080
```

### 方法2: ローカルで実行

#### Goのインストール

Goがインストールされていない場合は、以下のいずれかの方法でインストールしてください。

**Homebrewを使用（Mac）**

```bash
brew install go
```

**公式サイトからダウンロード**

https://go.dev/dl/

#### 依存関係のインストール

```bash
go mod tidy
```

#### 開発サーバーの起動

```bash
go run main.go
```

#### ビルドして実行

```bash
# ビルド
go build -o giter main.go

# 実行
./giter
```

サーバーが起動したら、ブラウザで以下のURLにアクセス:

```
http://localhost:8080
```

## 📂 プロジェクト構造

```
.
├── main.go                  # メインアプリケーション（Ginサーバー + GitHub API連携）
├── go.mod                   # Go依存関係管理
├── Dockerfile               # 本番環境用Dockerイメージ
├── Dockerfile.dev           # 開発環境用Dockerイメージ（ホットリロード対応）
├── docker-compose.yml       # 本番環境用Docker Compose設定
├── docker-compose.dev.yml   # 開発環境用Docker Compose設定
├── Makefile                 # 便利なコマンド集
├── .air.toml                # Airホットリロード設定
├── templates/
│   └── index.html           # フロントエンドHTML（Tailwind CSS + shadcn/ui）
├── static/                  # 静的ファイル用ディレクトリ
└── log/                     # ログファイル出力先（自動生成）
    └── YYYYMM/
        └── YYYYMMDD/
            └── app.log
```

## 🎨 機能

### Git履歴表示

- リポジトリ名
- コミットメッセージ
- コミット番号（短縮形）
- コミット時間
- GitHubへのリンク

### UI機能

- **カードレイアウト**: 各コミットをカードで表示
- **ソート**: 最新のコミットが上に表示
- **リフレッシュボタン**: 右下のFABボタンでデータを再取得
- **レスポンシブデザイン**: モバイル・デスクトップ両対応

## 🔧 技術スタック

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **API**: GitHub REST API v3
- **ログライブラリ**: zerolog

### フロントエンド
- **CSS**: Tailwind CSS（CDN版）
- **デザインシステム**: shadcn/ui インスパイア
- **JavaScript**: Vanilla JS（Fetch API使用）

### インフラ
- **コンテナ化**: Docker
- **オーケストレーション**: Docker Compose
- **ホットリロード**: Air（開発環境）

## 📊 ログ機能

アプリケーションは構造化ログを出力します。

### ログの保存先

ログは以下の階層構造で保存されます：

```
log/
└── YYYYMM/
    └── YYYYMMDD/
        └── app.log
```

例：`log/202602/20260214/app.log`

### ログレベル

環境変数 `LOG_LEVEL` で設定可能：
- `debug`: デバッグ情報を含む詳細なログ
- `info`: 通常の動作情報（デフォルト）
- `warn`: 警告メッセージ
- `error`: エラーメッセージのみ

docker-compose.ymlで設定例：
```yaml
environment:
  - TZ=Asia/Tokyo
  - LOG_LEVEL=debug
```

### ログの内容

- アプリケーションの起動/終了
- APIリクエストの処理状況
- GitHub API呼び出しの詳細
- エラー発生時の詳細情報
- リポジトリとコミットの取得状況

## 📝 API エンドポイント

### GET `/api/git-history`

develop-sudaユーザーのすべてのpublicリポジトリのコミット履歴を取得

**レスポンス例:**

```json
[
  {
    "repository_name": "example-repo",
    "commit_message": "Initial commit",
    "commit_sha": "a1b2c3d",
    "commit_time": "2024-01-01T12:00:00Z",
    "commit_url": "https://github.com/develop-suda/example-repo/commit/a1b2c3d4..."
  }
]
```

## 🎯 今後の拡張可能性

- ユーザー名の動的切り替え
- コミット数の統計表示
- 日付範囲フィルター
- リポジトリ別フィルター
- コミットメッセージの検索機能

## 📄 ライセンス

MIT License
