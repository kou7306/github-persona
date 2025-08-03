package funcs

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// 交互にトークンを取得するための関数
func GetTokens(currentIndex int) (string, int) {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return "", currentIndex
	}
	
	// 個人アクセストークンを環境変数から取得
	token := os.Getenv("GITHUB_TOKEN")
	
	if token == "" {
		fmt.Println("GITHUB_TOKEN not found in environment variables")
		return "", currentIndex
	}

	// ランダムシードの初期化
	rand.New(rand.NewSource(time.Now().UnixNano()))

	return token, currentIndex
}
