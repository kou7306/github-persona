package funcs

import (
	"fmt"
	"log"
)

func CreateUserStats(username string, star int) UserStats {
	// GitHubのアクセストークンを設定
	token, _ := GetTokens(0)
	stats := GetUserInfo(username, token)
	fmt.Println("stats: ", stats)
	ImgBytes, _ := GenerateGitHubStatsImage(stats, star, 600, 400)
	imageFileName := fmt.Sprintf("./images/stats_%s.png", username)

	err := SaveImage(imageFileName, ImgBytes)
	if err != nil {
		// エラーが発生した場合の処理
		fmt.Printf("Error saving stats image for %s: %v\n", username, err)
		log.Fatal(err)
	}
	
	fmt.Printf("Stats image saved successfully for %s: %s\n", username, imageFileName)

	return stats
}
