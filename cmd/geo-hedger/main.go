package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/nao-nba/geo-hedger/internal/asset"
	"github.com/nao-nba/geo-hedger/internal/io"
)

func main() {
	category := flag.String("type", "", "資産のカテゴリー (currency, western, eastern, stateless)")
	symbol := flag.String("symbol", "", "銘柄コード (例: USD, GLDM, BTC, 1929)")
	amount := flag.Float64("amount", 0, "数量")
	flag.Parse()

	filePath := "data/state.json"

	// 1. スプレッドシートから最新データ（レート & カテゴリー）を取得
	fmt.Println("--- 資産評価リレー開始 ---")
	symbolMap, err := io.FetchSymbolMapFromCSV()
	if err != nil {
		log.Fatal("❌ データの取得に失敗しました:", err)
	}

	// 2. 既存の資産データを読み込む
	state, err := io.LoadState(filePath)
	if err != nil {
		fmt.Println("INFO: 既存のデータがないため、新規作成します。")
		state = &asset.State{
			UpdatedAt: time.Now(),
			Assets: asset.Portfolio{
				Currencies:    make(map[string]float64),
				WesternRisk:   make(map[string]float64),
				EasternRisk:   make(map[string]float64),
				StatelessRisk: make(map[string]float64),
			},
		}
	}

	// 3. 引数がある場合は、スプレッドシートの定義と照合してバリデーション
	if *category != "" && *symbol != "" {
		info, exists := symbolMap[*symbol]

		// チェック①: スプレッドシートに存在するシンボルか？
		if !exists && *symbol != "JPY" {
			log.Fatalf("❌ エラー: シンボル '%s' はスプレッドシートに存在しません。", *symbol)
		}

		// チェック②: 指定されたカテゴリーは、スプレッドシートの定義と一致するか？
		if *symbol != "JPY" && info.Category != *category {
			log.Fatalf("❌ 属性エラー: '%s' は '%s' ではなく '%s' 属性です。", *symbol, *category, info.Category)
		}

		// チェックを通過したら登録
		switch *category {
		case "currency":
			state.Assets.Currencies[*symbol] = *amount
		case "western":
			state.Assets.WesternRisk[*symbol] = *amount
		case "eastern":
			state.Assets.EasternRisk[*symbol] = *amount
		case "stateless":
			state.Assets.StatelessRisk[*symbol] = *amount
		default:
			log.Fatal("❌ 不明なカテゴリーです。")
		}

		state.UpdatedAt = time.Now()

		if err := io.SaveState(filePath, state); err != nil {
			log.Fatal("❌ データの保存に失敗しました:", err)
		}
		fmt.Printf("✅ %s に %s: %.2f を登録・更新しました。\n", *category, *symbol, *amount)
	}

	// 4. 評価と計算（symbolMap からレートだけのマップを抽出して渡す）
	rates := make(asset.Rates)
	for sym, info := range symbolMap {
		rates[sym] = info.Rate
	}

	evaluated := state.Assets.Evaluate(rates)
	total := evaluated.CalculateTotal()

	fmt.Println("========================================")
	fmt.Printf("現在の総資産（円換算）: %.0f 円\n", total)
	fmt.Println("========================================")
}