package graphs

import (
	"fmt"
	"image/color"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// gonumライブラリを使用した理想的な積み上げ棒グラフ
func DrawCommitChart(commitsHistory []int, maxCommits int, width int, height int, username string) error {
	// 月ごとのデータを集計
	monthlyData := make(map[string]int)
	now := time.Now()
	
	for i, commits := range commitsHistory {
		// 過去の日付を計算
		pastDate := now.AddDate(0, 0, -(len(commitsHistory) - i - 1))
		monthKey := pastDate.Format("Jan") // "Jan"形式（年を削除）
		monthlyData[monthKey] += commits
	}

	// 時系列順にソート
	months := make([]string, 0, len(monthlyData))
	for month := range monthlyData {
		months = append(months, month)
	}
	sort.Strings(months)

	// グラフ用のデータを作成
	var values []float64
	var labels []string
	for _, month := range months {
		values = append(values, float64(monthlyData[month]))
		labels = append(labels, month)
	}

	// データが空の場合はエラーを返す
	if len(values) == 0 {
		return fmt.Errorf("no data to plot")
	}

	p := plot.New()

	// 背景色を設定
	bgColor := color.RGBA{R: 30, G: 30, B: 30, A: 255} // より暗い背景
	p.BackgroundColor = bgColor

	// 棒グラフを作成
	bars, err := plotter.NewBarChart(plotter.Values(values), vg.Points(50))
	if err != nil {
		return err
	}
	bars.Color = color.RGBA{R: 0, G: 180, B: 255, A: 255} // 明るい青
	bars.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白い枠線
	bars.LineStyle.Width = vg.Points(1)
	p.Add(bars)

	// タイトル
	p.Title.Text = "Monthly Commit Activity"
	p.Title.TextStyle.Font.Size = 45
	p.Title.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Title.Padding = 45

	// Y軸
	p.Y.Label.Text = "Commits"
	p.Y.Label.TextStyle.Font.Size = 36
	p.Y.Label.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Label.Font.Size = 30
	p.Y.Tick.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.LineStyle.Width = vg.Points(2)
	p.Y.Label.Padding = 45

	// Y軸の範囲を調整
	maxValue := 0.0
	for _, v := range values {
		if v > maxValue {
			maxValue = v
		}
	}
	if maxValue > 0 {
		p.Y.Max = maxValue * 1.2
	}

	// X軸
	p.X.Label.Text = "Month"
	p.X.Label.TextStyle.Font.Size = 36
	p.X.Label.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Label.Font.Size = 30
	p.X.Tick.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.LineStyle.Width = vg.Points(2)
	p.X.Label.Padding = 45

	// カスタムX軸ティック
	var ticks []plot.Tick
	for i, label := range labels {
		ticks = append(ticks, plot.Tick{Value: float64(i), Label: label})
	}
	p.X.Tick.Marker = plot.ConstantTicks(ticks)

	// グリッドを追加
	grid := plotter.NewGrid()
	grid.Horizontal.Color = color.RGBA{R: 80, G: 80, B: 80, A: 150}
	grid.Vertical.Color = color.RGBA{R: 80, G: 80, B: 80, A: 150}
	grid.Horizontal.Width = vg.Points(1.5)
	grid.Vertical.Width = vg.Points(1.5)
	p.Add(grid)

	p.X.Padding, p.Y.Padding = 0, 0
	imageFileName := fmt.Sprintf("./images/commits_history_%s.png", username)
	
	// 解像度を高く設定
	// 1インチ = 96ピクセルとして計算
	widthInch := vg.Inch * vg.Length(width) / 96
	heightInch := vg.Inch * vg.Length(height) / 96
	
	if err := p.Save(widthInch, heightInch, imageFileName); err != nil {
		return err
	}

	return nil
}
