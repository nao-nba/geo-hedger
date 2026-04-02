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

func FetchRatesFromCSV() (asset.Rates, error) {
	// 1. .envファイルを読み込む
	if err := godotenv.Load(); err != nil {
		fmt.Println("INFO: .env file not found, using system environment variables")
	}

	// 2. 環境変数からURLを取得
	csvURL := os.Getenv("GEO_HEDGER_CSV_URL")
	if csvURL == "" {
		return nil, fmt.Errorf("環境変数 GEO_HEDGER_CSV_URL が設定されていません")
	}

	rates := make(asset.Rates)

	// 3. HTTPリクエストでCSVデータを取得
	resp, err := http.Get(csvURL)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("サーバーがエラーを返しました: %d", resp.StatusCode)
	}

	// 4. CSVをパース
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSVのパースに失敗しました: %w", err)
	}

	// 5. 1行ずつ読み込んでマップに格納
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 2 {
			continue 
		}

		symbol := record[0]
		rateStr := record[1]

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			continue
		}

		rates[symbol] = rate
	}

	return rates, nil
}