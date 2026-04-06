// スプレッドシートから市場レートのCSVを取得する
package io

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nao-nba/geo-hedger/internal/asset"
)

// FetchSymbolMapFromCSV は環境変数のURLからレートとカテゴリーのマップを取得する
func FetchSymbolMapFromCSV() (asset.SymbolMap, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("INFO: .env file not found, using system environment variables")
	}

	csvURL := os.Getenv("GEO_HEDGER_CSV_URL")
	if csvURL == "" {
		return nil, fmt.Errorf("環境変数 GEO_HEDGER_CSV_URL が設定されていません")
	}

	symbolMap := make(asset.SymbolMap)

	resp, err := http.Get(csvURL)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("サーバーがエラーを返しました: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSVのパースに失敗しました: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue // ヘッダーをスキップ
		}
		// A列, B列, C列の3つが必要
		if len(record) < 3 {
			continue 
		}

		symbol := record[0]
		rateStr := record[1]
		category := record[2] // C列を取得

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			continue
		}

		// マップに構造体として格納
		symbolMap[symbol] = asset.SymbolInfo{
			Rate:     rate,
			Category: category,
		}
	}

	return symbolMap, nil
}