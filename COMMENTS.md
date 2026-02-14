# コードコメント追加について

このドキュメントでは、Giterアプリケーションに追加されたコメントの概要を説明します。

## コメント記法

### Go言語のコメント

**複数行コメント: `/* */`**
```go
/*
Repository はGitHub APIから取得するリポジトリ情報を表す構造体
GitHub API v3のリポジトリレスポンスの一部フィールドをマッピング
*/
type Repository struct {
    Name string `json:"name"` // リポジトリ名
}
```

**単一行コメント: `//`**
```go
Name string `json:"name"` // リポジトリ名（例: "my-project"）
```

### JavaScriptのコメント

**JSDocスタイル: `/** */`**
```javascript
/**
 * 関数の説明
 * @param {string} text - 引数の説明
 * @returns {string} 戻り値の説明
 */
function example(text) {
    return text;
}
```

**複数行コメント: `/* */`**
```javascript
/*
 * 処理の説明
 * 複数行にわたる詳細な説明
 */
```

**単一行コメント: `//`**
```javascript
const value = 10; // 変数の説明
```

## 追加されたコメントの種類

### 1. main.go（バックエンド）

#### 構造体のコメント
- `Repository`: GitHub APIから取得するリポジトリ情報の構造体
- `Commit`: GitHub APIから取得するコミット情報の構造体
- `CommitHistory`: フロントエンドに返却するレスポンス用の構造体
- 各フィールドに用途とデータ形式を明記

#### 定数のコメント
- `githubAPIBase`: GitHub REST API v3のベースURL
- `username`: 取得対象のGitHubユーザー名

#### 関数のコメント

**main関数**
- 目的: アプリケーションのエントリーポイント
- 処理内容:
  - Ginエンジンの初期化
  - CORSミドルウェアの設定と各パラメータの意味
  - 静的ファイルの配信設定
  - HTMLテンプレートの読み込み
  - ルーティング設定
  - サーバー起動

**getGitHistory関数**
- 目的: Git履歴を取得するAPIハンドラー
- 引数: `c *gin.Context` - Ginのコンテキスト
- 戻り値: HTTPレスポンス（JSON形式）
- 処理の流れ:
  1. リポジトリ一覧を取得
  2. 各リポジトリのコミット履歴を取得
  3. データを整形してJSON形式で返却
- エラーハンドリング: 個別リポジトリのエラーは全体の処理を停止しない

**fetchRepositories関数**
- 目的: GitHub APIから指定ユーザーの公開リポジトリ一覧を取得
- 戻り値: `[]Repository` と `error`
- API仕様へのリンク
- 注意事項:
  - レート制限（60リクエスト/時間）
  - 最大取得件数（100件）
- HTTPリクエストの各ステップを詳細に説明
- タイムアウト設定の理由
- リソース管理（defer resp.Body.Close()）

**fetchCommits関数**
- 目的: 指定されたリポジトリのコミット履歴を取得
- 引数: `repoFullName string` - リポジトリのフルネーム
- 戻り値: `[]Commit` と `error`
- API仕様へのリンク
- 注意事項:
  - デフォルトブランチのみ取得
  - 最大取得件数（100件）
- エラーケースの説明（404, 403など）

### 2. index.html（フロントエンド）

#### JavaScriptのコメント

**DOMContentLoadedイベントリスナー**
- イベントの発火タイミング
- 初期表示が速い理由

**loadCommits関数**
- JSDocスタイルのコメント
- 処理の流れを番号付きリストで明記
- 使用するDOM要素の説明
- async/awaitの説明
- Fetch APIの使用方法
- エラーハンドリング
- 配列のソート処理
- XSS対策（textContentの使用理由）

**createCommitCard関数**
- JSDocスタイルのコメント
- 引数の詳細（@param タグ）
- 戻り値の説明（@returns タグ）
- 日時フォーマットの説明
- テンプレートリテラルの使用方法
- HTML内にもコメント追加:
  - 各UI要素の役割
  - セキュリティ対策（rel="noopener noreferrer"）

**escapeHtml関数**
- JSDocスタイルのコメント
- XSS攻撃の説明
- 変換対象文字の一覧と理由
- 具体例（入力と出力）
- 正規表現の説明
- マップの使用方法

## コメントの記述方針

### 1. 複数行コメントは `/* */` 形式を使用

**Go言語:**
```go
/*
関数の説明
複数行にわたる詳細な説明
*/
func example() {
    /* 処理ブロックの説明 */
    doSomething()
}
```

**JavaScript:**
```javascript
/**
 * JSDoc形式（関数のドキュメント）
 * @param {type} name - 説明
 */
function example(name) {
    /* 処理ブロックの説明 */
    doSomething();
}
```

### 2. 「何をしているか」だけでなく「なぜそうするのか」を説明
例:
```go
/*
	deferでレスポンスボディを確実にクローズ
	これによりリソースリークを防ぐ
*/
defer resp.Body.Close()
```

### 3. 引数・戻り値の詳細な説明
例:
```go
/*
引数:
  repoFullName string - リポジトリのフルネーム（例: "develop-suda/project-name"）
                       所有者名とリポジトリ名をスラッシュで結合した形式

戻り値:
  []Commit - 取得したコミット情報のスライス
  error - エラーが発生した場合のエラーオブジェクト
*/
```

### 4. セキュリティやパフォーマンスに関する注意事項
例:
```go
/*
注意:
  - GitHub APIは認証なしで60リクエスト/時間の制限あり
  - タイムアウトは10秒に設定（長時間のリクエストを防ぐ）
*/
```

### 5. エラーケースの説明
例:
```go
/*
	HTTPステータスコードが200 OK以外の場合はエラー
	404 Not Foundの場合はリポジトリが存在しないか、アクセス権限がない
	403 Forbiddenの場合はAPIレート制限に到達した可能性がある
*/
```

### 6. 主要な変数の意図
例:
```go
/*
	allCommitsは全リポジトリのコミット履歴を格納するスライス
	初期容量は指定せず、append()で動的に拡張
*/
var allCommits []CommitHistory
```

## コメントの効果

### 可読性の向上
- 新しい開発者が素早くコードを理解できる
- メンテナンス時の理解コストを削減

### ドキュメントとしての役割
- 外部ドキュメントを参照せずにコード内で完結
- APIのリンクなど、関連情報へのアクセスが容易

### セキュリティ意識の向上
- XSS対策などのセキュリティ上の配慮が明確
- なぜその対策が必要なのかを理解できる

### デバッグの効率化
- 各処理の意図が明確なため、問題箇所の特定が容易
- エラーケースの説明により、トラブルシューティングが簡単

## 今後の保守について

コードを変更する際は、以下のルールに従ってください：

### 1. 新しい関数を追加する場合

**Go言語:**
```go
/*
関数名 は処理の説明

引数:
  arg1 type - 引数の説明
  arg2 type - 引数の説明

戻り値:
  type - 戻り値の説明
  error - エラーの説明

注意:
  - 特記事項がある場合
*/
func funcName(arg1 type, arg2 type) (type, error) {
    /* 処理内容 */
}
```

**JavaScript:**
```javascript
/**
 * 関数の説明
 *
 * @param {type} arg1 - 引数の説明
 * @param {type} arg2 - 引数の説明
 * @returns {type} 戻り値の説明
 */
function funcName(arg1, arg2) {
    /* 処理内容 */
}
```

### 2. 既存のコードを変更する場合
- 該当箇所のコメントも更新
- 変更理由を追記（必要に応じて）
- `/* */` 形式を維持

### 3. 複雑なロジックを実装する場合
```go
/*
	なぜその実装方法を選んだのかを説明

	代替案: 他の実装方法の説明
	理由: なぜこの方法を選んだのか
*/
```

## 参考

- [GitHub REST API ドキュメント](https://docs.github.com/ja/rest)
- [Goの効果的なコメントの書き方](https://go.dev/doc/effective_go#commentary)
- [JSDocスタイルガイド](https://jsdoc.app/)
