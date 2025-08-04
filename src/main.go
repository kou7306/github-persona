package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tomoish/github-persona/funcs"
	"github.com/tomoish/github-persona/graphs"
)

// func handler(w http.ResponseWriter, r *http.Request) {

// 	ctx := context.Background()

// 	path := r.URL.Path
// 	segments := strings.Split(path, "/")
// 	username := segments[1]

// 	client := github.NewClient(nil)
// 	repos, _, _ := client.Repositories.ListByUser(ctx, username, nil)
// 	for _, repo := range repos {
// 		repoName := *repo.Name
// 		stars := *repo.StargazersCount
// 		forks := *repo.ForksCount

// 		fmt.Fprintf(w, "Repo: %v, Stars: %d, Forks: %d\n", repoName, stars, forks)
// 	}
// 	fmt.Fprint(w, repos)
// }

// // 言語画像生成
// func getLanguageHandler(w http.ResponseWriter, r *http.Request) {
// 	funcs.CreateLanguageImg()
// }

// //キャラ画像生成

// func getCharacterHandler(w http.ResponseWriter, r *http.Request) {
// 	funcs.CreateCharacterImg()
// }

// // 全て合体
// func mergeAllContents(w http.ResponseWriter, r *http.Request) {
// 	funcs.Merge_all("./images/background.png", "./images/stats.png", "./images/generate_character.png", "./images/language.png", "./images/commits_history.png")
// }

// // 背景生成
// func getBackgroundHandler(w http.ResponseWriter, r *http.Request) {
// 	funcs.DrawBackground("Lv.30", "神")
// }

// func getCommitStreakHandler(w http.ResponseWriter, r *http.Request) {

// 	queryValues := r.URL.Query()
// 	username := queryValues.Get("username")

// 	if username == "" {
// 		http.Error(w, "username is required", http.StatusBadRequest)
// 		return
// 	}

// 	streak, dailyCommits, _, err := funcs.GetCommitHistory(username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprint(w, streak, dailyCommits)

// }

// func getHistoryHandler(w http.ResponseWriter, r *http.Request) {

// 	queryValues := r.URL.Query()
// 	username := queryValues.Get("username")

// 	if username == "" {
// 		http.Error(w, "username is required", http.StatusBadRequest)
// 		return
// 	}

// 	_, dailyCommits, maxCommits, err := funcs.GetCommitHistory(username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = graphs.DrawCommitChart(dailyCommits, maxCommits, 1000, 700, username)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	http.ServeFile(w, r, "./images/commits_history.png")
// }
// func getuserHandler(w http.ResponseWriter, r *http.Request) {

// 	username := "kou7306"
// 	// GitHubのアクセストークンを設定
// 	token, _ := funcs.GetTokens(0)
// 	stats := funcs.GetUserInfo(username, token)
// 	fmt.Println("stats: ", stats)
// 	ImgBytes, _ := funcs.GenerateGitHubStatsImage(stats, 600, 400)
// 	fmt.Println("ImgBytes: ", ImgBytes)

// 	err := funcs.SaveImage("images/stats.png", ImgBytes)
// 	if err != nil {
// 		// エラーが発生した場合の処理
// 		log.Fatal(err) // または他のエラーハンドリング方法を選択してください
// 	}

// }

// 画像生成エンドポイント

func createHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                   // すべてのオリジンからのアクセスを許可
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // 許可するHTTPメソッド
	w.Header().Set("Access-Control-Allow-Headers", "*") // すべてのヘッダーを許可

	// OPTIONSリクエストへの対応（プリフライトリクエスト）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	queryValues := r.URL.Query()
	username := queryValues.Get("username")
	if r.Method == http.MethodGet {
		// GETリクエストの処理
		fmt.Printf("Processing request for user: %s\n", username)
		
		// キャッシュチェック - 既存の画像があるか確認
		nocache := queryValues.Get("nocache")
		finalImagePath := fmt.Sprintf("./images/final_%s.png", username)
		if nocache != "1" {
			if _, err := os.Stat(finalImagePath); err == nil {
				// キャッシュされた画像が存在する場合、直接返す
				fmt.Printf("Cache hit for %s, serving cached image\n", username)
				imageBytes, err := os.ReadFile(finalImagePath)
				if err != nil {
					fmt.Printf("Error reading cached image for %s: %v\n", username, err)
					http.Error(w, "Failed to read cached image", http.StatusInternalServerError)
					return
				}
				// レスポンスヘッダーを設定
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Cache-Control", "public, max-age=3600")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
				w.Write(imageBytes)
				fmt.Printf("Cached image served successfully for %s\n", username)
				return
			}
		}
		fmt.Printf("Cache miss for %s, generating new image\n", username)
		
		// 並行処理で画像生成を高速化
		reposChan := make(chan []funcs.Repository)
		statsChan := make(chan funcs.UserStats)
		languageChan := make(chan []funcs.LanguageStat)
		commitChan := make(chan []int)
		maxCommitsChan := make(chan int)
		errChan := make(chan error)
		
		// リポジトリ情報を並行取得（一度だけ）
		go func() {
			repos, star, _ := funcs.GetRepositories(username)
			reposChan <- repos
			// 同じリポジトリ情報を使って統計も生成
			stats := funcs.CreateUserStats(username, star)
			statsChan <- stats
		}()
		
		// 言語情報を並行取得
		go func() {
			language := funcs.CreateLanguageImg(username)
			languageChan <- language
		}()
		
		// コミット履歴を並行取得
		go func() {
			_, dailyCommits, maxCommits, err := funcs.GetCommitHistory(username)
			if err != nil {
				errChan <- err
				return
			}
			commitChan <- dailyCommits
			maxCommitsChan <- maxCommits
		}()
		
		// 結果を待機
		_ = <-reposChan // 未使用変数を削除
		stats := <-statsChan
		language := <-languageChan
		dailyCommits := <-commitChan
		maxCommits := <-maxCommitsChan
		
		// エラーチェック
		select {
		case err := <-errChan:
			fmt.Printf("Error getting commit history for %s: %v\n", username, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		default:
		}
		
		fmt.Printf("All data retrieved for %s\n", username)
		
		total := stats.TotalStars + stats.ContributedTo + stats.TotalIssues + stats.TotalPRs + stats.TotalCommits
		
		//レベル、職業判定
		profession, level := funcs.JudgeRank(language, stats, 0) // starの代わりに0を使用
		fmt.Printf("Rank judged for %s: profession=%s, level=%d\n", username, profession, level)
		
		//対象のキャラの画像を取得
		img := funcs.DispatchPictureBasedOnProfession(profession)
		fmt.Printf("Character image dispatched for %s: %s\n", username, img)

		filePath := fmt.Sprintf("characterImages/%s", img)

		// 並行処理で画像生成（すべて同時に実行）
		backgroundChan := make(chan error)
		characterChan := make(chan error)
		commitChartChan := make(chan error)
		statsImgChan := make(chan error)
		languageImgChan := make(chan error)
		
		// 背景画像を並行生成
		go func() {
			funcs.DrawBackground(username, "Lv."+strconv.Itoa(level), profession)
			backgroundChan <- nil
		}()
		
		// キャラクター画像を並行生成
		go func() {
			funcs.CreateCharacterImg(filePath, "images/gauge.png", total, level, username)
			characterChan <- nil
		}()
		
		// コミットチャートを並行生成（解像度を下げて高速化）
		go func() {
			err := graphs.DrawCommitChart(dailyCommits, maxCommits, 800, 500, username) // 解像度を下げて高速化
			commitChartChan <- err
		}()
		
		// 統計画像の生成も並行化（CreateUserStats内で生成される）
		go func() {
			// 統計画像を生成
			funcs.CreateUserStats(username, stats.TotalStars)
			statsImgChan <- nil
		}()
		
		// 言語画像の生成も並行化（CreateLanguageImg内で生成される）
		go func() {
			// CreateLanguageImgは既に実行済みなので、ここでは待機のみ
			languageImgChan <- nil
		}()
		
		// すべての画像生成完了を待機
		<-backgroundChan
		<-characterChan
		<-statsImgChan
		<-languageImgChan
		if err := <-commitChartChan; err != nil {
			fmt.Printf("Error drawing commit chart for %s: %v\n", username, err)
			http.Error(w, "Failed to draw commit chart", http.StatusInternalServerError)
			return
		}
		
		fmt.Printf("All images generated for %s\n", username)
		
		backImg := fmt.Sprintf("./images/background_%s.png", username)
		statsImg := fmt.Sprintf("./images/stats_%s.png", username)
		characterImg := fmt.Sprintf("./images/generate_character_%s.png", username)
		languageImg := fmt.Sprintf("./images/language_%s.png", username)
		dateImg := fmt.Sprintf("./images/commits_history_%s.png", username)
		
		// 画像ファイルの存在確認
		requiredFiles := []string{backImg, statsImg, characterImg, languageImg, dateImg}
		for _, file := range requiredFiles {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				fmt.Printf("Error: Required file %s does not exist\n", file)
				http.Error(w, fmt.Sprintf("Required file %s not found", file), http.StatusInternalServerError)
				return
			}
		}
		
		// 全て合体して画像をメモリ上で生成
		imageBytes, err := funcs.Merge_all_to_bytes(backImg, statsImg, characterImg, languageImg, dateImg)
		if err != nil {
			fmt.Printf("Error merging images for %s: %v\n", username, err)
			http.Error(w, "Failed to generate image", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Images merged successfully for %s\n", username)

		// 最終的なマージされた画像をファイルとして保存
		finalImagePath = fmt.Sprintf("./images/final_%s.png", username)
		absPath, _ := os.Getwd()
		fullPath := fmt.Sprintf("%s/%s", absPath, finalImagePath)
		err = os.WriteFile(finalImagePath, imageBytes, 0644)
		if err != nil {
			fmt.Printf("Error saving final image for %s: %v\n", username, err)
			http.Error(w, "Failed to save final image", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Final image saved successfully for %s\n", username)
		fmt.Printf("  Path: %s\n", fullPath)
		fmt.Printf("  Size: %d bytes\n", len(imageBytes))

		// 部品画像のみを削除（並行処理）
		go func() {
			gaugeImg := "images/gauge.png"
			os.Remove(backImg)
			os.Remove(statsImg)
			os.Remove(characterImg)
			os.Remove(languageImg)
			os.Remove(dateImg)
			os.Remove(gaugeImg)
			fmt.Printf("Component images cleaned up for %s\n", username)
		}()

		// レスポンスヘッダーを設定
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1時間キャッシュ
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		// 画像データを直接レスポンスとして返す
		w.Write(imageBytes)
		fmt.Printf("Response sent successfully for %s\n", username)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GitHub用の最適化されたエンドポイント
func githubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	queryValues := r.URL.Query()
	username := queryValues.Get("username")
	
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// キャッシュバスティングパラメータをチェック
	forceUpdate := queryValues.Get("update") == "1"
	nocache := queryValues.Get("nocache") == "1"
	timestamp := queryValues.Get("t") != "" // タイムスタンプパラメータがある場合は強制更新
	
	// キャッシュチェック（強制更新でない場合のみ）
	if !forceUpdate && !nocache && !timestamp {
		cachedImagePath := fmt.Sprintf("./images/final_%s.png", username)
		if _, err := os.Stat(cachedImagePath); err == nil {
			// キャッシュされた画像が存在する場合、直接返す
			imageBytes, readErr := os.ReadFile(cachedImagePath)
			if readErr == nil {
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Cache-Control", "public, max-age=3600") // 1時間キャッシュ（短縮）
				w.Write(imageBytes)
				return
			}
		}
	}

	// 通常の処理を実行
	createHandler(w, r)
}

func main() {
	// http.HandleFunc("/test", handler)
	// http.HandleFunc("/streak", getCommitStreakHandler)
	// http.HandleFunc("/language", getLanguageHandler)
	// http.HandleFunc("/character", getCharacterHandler)
	// http.HandleFunc("/history", getHistoryHandler)
	// http.HandleFunc("/user", getuserHandler)
	// http.HandleFunc("/merge", mergeAllContents)
	// http.HandleFunc("/background", getBackgroundHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/github", githubHandler)
	fmt.Println("Hello, World!")
	
	// タイムアウト設定付きのサーバー
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
