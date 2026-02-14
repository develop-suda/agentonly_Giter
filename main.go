package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*
Repository はGitHub APIから取得するリポジトリ情報を表す構造体
GitHub API v3のリポジトリレスポンスの一部フィールドをマッピング
*/
type Repository struct {
	Name        string `json:"name"`        // リポジトリ名（例: "my-project"）
	FullName    string `json:"full_name"`   // フルネーム（例: "develop-suda/my-project"）
	Description string `json:"description"` // リポジトリの説明文
	HTMLURL     string `json:"html_url"`    // GitHubのリポジトリURL
}

/*
Commit はGitHub APIから取得するコミット情報を表す構造体
GitHub API v3のコミットレスポンスの必要なフィールドをマッピング
*/
type Commit struct {
	SHA string `json:"sha"` // コミットハッシュ（40文字の16進数文字列）
	/* Commitフィールドはネストされた構造を持つ */
	Commit struct {
		Message string `json:"message"` // コミットメッセージ
		/* Authorフィールドはコミット作成者の情報 */
		Author struct {
			Name  string    `json:"name"`  // 作成者名
			Email string    `json:"email"` // 作成者のメールアドレス
			Date  time.Time `json:"date"`  // コミット作成日時（ISO 8601形式）
		} `json:"author"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"` // GitHubのコミットURL
}

/*
CommitHistory はフロントエンドに返却するレスポンス用の構造体
GitHub APIのレスポンスを整形し、必要な情報のみを含む
*/
type CommitHistory struct {
	RepositoryName string    `json:"repository_name"` // リポジトリ名
	CommitMessage  string    `json:"commit_message"`  // コミットメッセージ
	CommitSHA      string    `json:"commit_sha"`      // コミットハッシュ（短縮形、7文字）
	CommitTime     time.Time `json:"commit_time"`     // コミット作成日時
	CommitURL      string    `json:"commit_url"`      // GitHubのコミットページへのリンク
}

const (
	/* githubAPIBase はGitHub REST API v3のベースURL */
	githubAPIBase = "https://api.github.com"
	/*
		username は取得対象のGitHubユーザー名
		定数として定義することで、変更が容易になる
	*/
	username = "develop-suda"
)

/*
main はアプリケーションのエントリーポイント
Ginフレームワークを使用してWebサーバーを起動し、以下の機能を提供する：
- CORS対応（クロスオリジンリクエストを許可）
- 静的ファイルの配信
- HTMLテンプレートのレンダリング
- REST APIエンドポイント
*/
func main() {
	/*
		gin.Default()はロガーとリカバリーミドルウェアが組み込まれたGinエンジンを作成
		リカバリーミドルウェアはpanicを検知し、500エラーを返す
	*/
	r := gin.Default()

	/*
		CORS（Cross-Origin Resource Sharing）ミドルウェアの設定
		フロントエンドが異なるオリジンから API を呼び出せるようにする
	*/
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                      // すべてのオリジンからのアクセスを許可（本番環境では制限を推奨）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 許可するHTTPメソッド
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},        // 許可するリクエストヘッダー
		ExposeHeaders:    []string{"Content-Length"},                          // フロントエンドに公開するレスポンスヘッダー
		AllowCredentials: true,                                                // クッキーなどの認証情報の送信を許可
		MaxAge:           12 * time.Hour,                                      // プリフライトリクエストのキャッシュ時間
	}))

	/*
		静的ファイルの配信設定
		URLパス "/static" へのアクセスを "./static" ディレクトリにマッピング
		例: /static/css/style.css -> ./static/css/style.css
	*/
	r.Static("/static", "./static")

	/*
		HTMLテンプレートファイルの読み込み
		"templates/*" パターンに一致するすべてのファイルをテンプレートとして登録
	*/
	r.LoadHTMLGlob("templates/*")

	/*
		ルートページ（"/"）へのGETリクエストのハンドラー
		index.htmlテンプレートをレンダリングして返す
	*/
	r.GET("/", func(c *gin.Context) {
		/*
			第一引数: HTTPステータスコード（200 OK）
			第二引数: テンプレート名
			第三引数: テンプレートに渡すデータ（今回はnil）
		*/
		c.HTML(http.StatusOK, "index.html", nil)
	})

	/*
		Git履歴APIエンドポイント
		"/api/git-history" へのGETリクエストをgetGitHistory関数で処理
		このエンドポイントは全リポジトリのコミット履歴をJSON形式で返す
	*/
	r.GET("/api/git-history", getGitHistory)

	/* サーバー起動メッセージ */
	fmt.Println("Server starting on :8080")

	/*
		Webサーバーを起動し、ポート8080でリクエストを待ち受ける
		この関数はブロッキングで、サーバーが停止するまで戻らない
	*/
	r.Run(":8080")
}

/*
getGitHistory はGit履歴を取得するAPIハンドラー
処理の流れ:
1. fetchRepositories()でユーザーの全公開リポジトリを取得
2. 各リポジトリのコミット履歴をfetchCommits()で取得
3. 全コミットを統合してJSON形式で返却

引数:
  c *gin.Context - Ginのコンテキスト。リクエスト・レスポンス情報を含む

レスポンス:
  成功時: 200 OK, []CommitHistory（全コミット履歴のJSON配列）
  失敗時: 500 Internal Server Error, {"error": "エラーメッセージ"}
*/
func getGitHistory(c *gin.Context) {
	/*
		fetchRepositories()を呼び出し、対象ユーザーの全公開リポジトリを取得
		戻り値: repos（リポジトリのスライス）, err（エラー）
	*/
	repos, err := fetchRepositories()
	if err != nil {
		/*
			エラーが発生した場合、500エラーとエラーメッセージをJSON形式で返す
			gin.Hは map[string]interface{} のエイリアスで、JSON生成に使用
		*/
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*
		allCommitsは全リポジトリのコミット履歴を格納するスライス
		初期容量は指定せず、append()で動的に拡張
	*/
	var allCommits []CommitHistory

	/*
		各リポジトリをイテレートしてコミット履歴を取得
		range repos は (index, value) を返すが、indexは使用しないので _ で無視
	*/
	for _, repo := range repos {
		/* repo.FullName（例: "develop-suda/project-name"）を使用してコミットを取得 */
		commits, err := fetchCommits(repo.FullName)
		if err != nil {
			/*
				個別リポジトリのエラーは全体の処理を停止せず、ログ出力のみ
				これにより一部のリポジトリが取得できなくても他のリポジトリは表示される
			*/
			fmt.Printf("Error fetching commits for %s: %v\n", repo.Name, err)
			continue /* 次のリポジトリの処理に進む */
		}

		/* 取得したコミットをCommitHistory形式に変換してスライスに追加 */
		for _, commit := range commits {
			allCommits = append(allCommits, CommitHistory{
				RepositoryName: repo.Name,             // リポジトリ名
				CommitMessage:  commit.Commit.Message, // コミットメッセージ
				CommitSHA:      commit.SHA[:7],        // コミットハッシュを7文字に短縮（Gitの慣習）
				CommitTime:     commit.Commit.Author.Date, // コミット作成日時
				CommitURL:      commit.HTMLURL,        // GitHubのコミットページURL
			})
		}
	}

	/*
		全コミット履歴をJSON形式でレスポンスとして返す
		Ginが自動的にContent-Type: application/jsonヘッダーを設定
	*/
	c.JSON(http.StatusOK, allCommits)
}

/*
fetchRepositories はGitHub APIから指定ユーザーの公開リポジトリ一覧を取得する
GitHub REST API v3のリポジトリ一覧取得エンドポイントを使用
API仕様: https://docs.github.com/ja/rest/repos/repos#list-repositories-for-a-user

戻り値:
  []Repository - 取得したリポジトリ情報のスライス（最大100件）
  error - エラーが発生した場合のエラーオブジェクト、正常時はnil

注意:
  - GitHub APIは認証なしで60リクエスト/時間の制限あり
  - per_page=100で最大100件を取得（デフォルトは30件）
*/
func fetchRepositories() ([]Repository, error) {
	/*
		GitHub API URLを構築
		クエリパラメータ:
		  - type=public: 公開リポジトリのみ取得
		  - per_page=100: 1ページあたり100件（APIの最大値）
	*/
	url := fmt.Sprintf("%s/users/%s/repos?type=public&per_page=100", githubAPIBase, username)

	/*
		HTTPリクエストを作成
		第一引数: HTTPメソッド（GET）
		第二引数: リクエストURL
		第三引数: リクエストボディ（GETなのでnil）
	*/
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		/* リクエスト作成に失敗した場合（通常は発生しない） */
		return nil, err
	}

	/*
		GitHub API v3用のAcceptヘッダーを設定
		これによりAPI v3のレスポンス形式が保証される
	*/
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	/*
		HTTPクライアントを作成
		Timeout: 10秒でタイムアウト（長時間のリクエストを防ぐ）
	*/
	client := &http.Client{Timeout: 10 * time.Second}

	/* HTTPリクエストを実行 */
	resp, err := client.Do(req)
	if err != nil {
		/* ネットワークエラーやタイムアウトの場合 */
		return nil, err
	}
	/*
		deferでレスポンスボディを確実にクローズ
		これによりリソースリークを防ぐ
	*/
	defer resp.Body.Close()

	/* HTTPステータスコードが200 OK以外の場合はエラー */
	if resp.StatusCode != http.StatusOK {
		/* エラー詳細をレスポンスボディから読み取る */
		body, _ := io.ReadAll(resp.Body)
		/* ステータスコードとボディを含むエラーメッセージを返す */
		return nil, fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	/* レスポンスボディをRepository構造体のスライスにデコード */
	var repos []Repository
	/*
		json.NewDecoder()はストリーミング処理を行うため、大きなレスポンスでもメモリ効率が良い
	*/
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		/* JSONパースエラー（APIレスポンス形式が期待と異なる場合） */
		return nil, err
	}

	/* 取得したリポジトリ一覧を返す */
	return repos, nil
}

/*
fetchCommits は指定されたリポジトリのコミット履歴を取得する
GitHub REST API v3のコミット一覧取得エンドポイントを使用
API仕様: https://docs.github.com/ja/rest/commits/commits#list-commits

引数:
  repoFullName string - リポジトリのフルネーム（例: "develop-suda/project-name"）
                       所有者名とリポジトリ名をスラッシュで結合した形式

戻り値:
  []Commit - 取得したコミット情報のスライス（最大100件、新しい順）
  error - エラーが発生した場合のエラーオブジェクト、正常時はnil

注意:
  - デフォルトブランチのコミットのみ取得される
  - per_page=100で最大100件を取得（APIの最大値）
  - GitHub APIは認証なしで60リクエスト/時間の制限あり
*/
func fetchCommits(repoFullName string) ([]Commit, error) {
	/*
		GitHub API URLを構築
		エンドポイント: /repos/{owner}/{repo}/commits
		クエリパラメータ:
		  - per_page=100: 1ページあたり100件（APIの最大値）
	*/
	url := fmt.Sprintf("%s/repos/%s/commits?per_page=100", githubAPIBase, repoFullName)

	/*
		HTTPリクエストを作成
		第一引数: HTTPメソッド（GET）
		第二引数: リクエストURL
		第三引数: リクエストボディ（GETなのでnil）
	*/
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		/* リクエスト作成に失敗した場合（通常は発生しない） */
		return nil, err
	}

	/*
		GitHub API v3用のAcceptヘッダーを設定
		これによりAPI v3のレスポンス形式が保証される
	*/
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	/*
		HTTPクライアントを作成
		Timeout: 10秒でタイムアウト（長時間のリクエストを防ぐ）
	*/
	client := &http.Client{Timeout: 10 * time.Second}

	/* HTTPリクエストを実行 */
	resp, err := client.Do(req)
	if err != nil {
		/* ネットワークエラーやタイムアウトの場合 */
		return nil, err
	}
	/*
		deferでレスポンスボディを確実にクローズ
		これによりリソースリークを防ぐ
	*/
	defer resp.Body.Close()

	/*
		HTTPステータスコードが200 OK以外の場合はエラー
		404 Not Foundの場合はリポジトリが存在しないか、アクセス権限がない
		403 Forbiddenの場合はAPIレート制限に到達した可能性がある
	*/
	if resp.StatusCode != http.StatusOK {
		/* エラー詳細をレスポンスボディから読み取る */
		body, _ := io.ReadAll(resp.Body)
		/* ステータスコードとボディを含むエラーメッセージを返す */
		return nil, fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	/* レスポンスボディをCommit構造体のスライスにデコード */
	var commits []Commit
	/*
		json.NewDecoder()はストリーミング処理を行うため、大きなレスポンスでもメモリ効率が良い
	*/
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		/* JSONパースエラー（APIレスポンス形式が期待と異なる場合） */
		return nil, err
	}

	/* 取得したコミット一覧を返す（新しい順にソート済み） */
	return commits, nil
}
