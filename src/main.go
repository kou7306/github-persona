package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	// OPTIONSリクエストへの対応（プリフライトリクエスト）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	queryValues := r.URL.Query()
	username := queryValues.Get("username")
	if r.Method == http.MethodGet {
		// GETリクエストの処理
		// 画像生成の処理...
		fmt.Printf("Processing request for user: %s\n", username)
		
		_, star, _ := funcs.GetRepositories(username)
		fmt.Printf("Repositories processed for %s\n", username)
		
		// stats取得と画像生成
		stats := funcs.CreateUserStats(username, star)
		fmt.Printf("Stats created for %s\n", username)
		
		total := stats.TotalStars + stats.ContributedTo + stats.TotalIssues + stats.TotalPRs + stats.TotalCommits
		
		// 言語画像の生成
		language := funcs.CreateLanguageImg(username)
		fmt.Printf("Language image created for %s\n", username)
		
		//レベル、職業判定
		profession, level := funcs.JudgeRank(language, stats, star)
		fmt.Printf("Rank judged for %s: profession=%s, level=%d\n", username, profession, level)
		
		//対象のキャラの画像を取得
		img := funcs.DispatchPictureBasedOnProfession(profession)
		fmt.Printf("Character image dispatched for %s: %s\n", username, img)

		filePath := fmt.Sprintf("characterImages/%s", img)

		// 背景画像の生成
		funcs.DrawBackground(username, "Lv."+strconv.Itoa(level), profession)
		fmt.Printf("Background drawn for %s\n", username)

		// キャラクター画像の生成
		funcs.CreateCharacterImg(filePath, "images/gauge.png", total, level, username)
		fmt.Printf("Character image created for %s\n", username)

		_, dailyCommits, maxCommits, err := funcs.GetCommitHistory(username)
		if err != nil {
			fmt.Printf("Error getting commit history for %s: %v\n", username, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Commit history retrieved for %s\n", username)

		err = graphs.DrawCommitChart(dailyCommits, maxCommits, 1000, 700, username)
		if err != nil {
			fmt.Printf("Error drawing commit chart for %s: %v\n", username, err)
			http.Error(w, "Failed to draw commit chart", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Commit chart drawn for %s\n", username)
		
		backImg := fmt.Sprintf("./images/background_%s.png", username)
		statsImg := fmt.Sprintf("./images/stats_%s.png", username)
		characterImg := fmt.Sprintf("./images/generate_character_%s.png", username)
		languageImg := fmt.Sprintf("./images/language_%s.png", username)
		dateImg := fmt.Sprintf("./images/commits_history_%s.png", username)
		
		// 全て合体して画像をメモリ上で生成
		imageBytes, err := funcs.Merge_all_to_bytes(backImg, statsImg, characterImg, languageImg, dateImg)
		if err != nil {
			fmt.Printf("Error merging images for %s: %v\n", username, err)
			http.Error(w, "Failed to generate image", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Images merged successfully for %s\n", username)

		// レスポンスヘッダーを設定
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// 画像データを直接レスポンスとして返す
		w.Write(imageBytes)
		fmt.Printf("Response sent successfully for %s\n", username)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

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
	fmt.Println("Hello, World!")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
