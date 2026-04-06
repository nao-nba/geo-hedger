package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/nao-nba/geo-hedger/internal/asset"
	"github.com/nao-nba/geo-hedger/internal/io"
)

// 多言語化のための辞書
var i18n = map[string]map[string]string{
	"ja": {
		"relay_start":     "--- 資産評価リレー開始 ---",
		"target_phase":   "【現在の目標フェーズ】",
		"scenario":       "【シナリオ】",
		"action":         "【行動方針】",
		"total_asset":    "現在の総資産（余剰資金）",
		"diff_title":     "■ 資産配分と目標フェーズとの差分",
		"table_header":   "[シンボル] | 現在の比率 | 目標比率 | 乖離額（目標 - 現在）",
		"status_opt":     "適正",
		"status_short":   "不足",
		"status_excess":  "超過",
		"next_actions":   "■ NEXT ACTIONS（誰が、どこで、何をすべきか）",
		"no_rebalance":   "● 現在、リバランスの必要はありません。良好な分散が保たれています。",
		"sell":           "【売却・円転】",
		"buy":            "【購入】",
		"none":           "● 特になし",
		"default_owner":  "推奨口座（任意）",
		"sell_action":    "● %s: %s を %.0f 円分 売却/円転\n",
		"buy_action":     "● %s: %s を %.0f 円分 購入\n",
		
	},
	"en": {
		"relay_start":     "--- Asset Evaluation Relay Started ---",
		"target_phase":   "[Current Target Phase]",
		"scenario":       "[Scenario]",
		"action":         "[Action Plan]",
		"total_asset":    "Current Total Assets (Surplus Funds)",
		"diff_title":     "■ Asset Allocation & Gap from Target Phase",
		"table_header":   "[Symbol]   | Current %  | Target %  | Gap (Target - Current)",
		"status_opt":     "Optimal",
		"status_short":   "Short",
		"status_excess":  "Excess",
		"next_actions":   "■ NEXT ACTIONS (Who, Where, and What to do)",
		"no_rebalance":   "● No rebalancing needed at this time.",
		"sell":           "[SELL / DIVEST]",
		"buy":            "[BUY / ACCUMULATE]",
		"none":           "● None",
		"default_owner":  "Recommended Account",
		"sell_action":    "● %s: SELL %s worth of %.0f JPY\n",
		"buy_action":     "● %s: BUY %s worth of %.0f JPY\n",
	},
}

func main() {
	owner := flag.String("owner", "", "資産の持ち主 (例: 家族A（SBI証券）)")
	symbol := flag.String("symbol", "", "銘柄コード (例: USD, GLDM, BTC, 1929)")
	amount := flag.Float64("amount", 0, "数量")
	
	// ★ 言語の引数（デフォルトは英語）
	lang := flag.String("lang", "en", "Language / 表示言語 (en, ja)")
	
	flag.Parse()

	// 指定された言語が辞書になければ、強制的に日本語にする
	if _, ok := i18n[*lang]; !ok {
		*lang = "ja"
	}

	filePath := "data/state.json"

	// 1. スプレッドシートから最新データ（レート & カテゴリー）を取得
	fmt.Println(i18n[*lang]["relay_start"])
	symbolMap, err := io.FetchSymbolMapFromCSV()
	if err != nil {
		log.Fatal("❌ データの取得に失敗しました:", err)
	}

	// 2. 既存の資産データを読み込む
	state, err := io.LoadState(filePath)
	if err != nil {
		fmt.Println("INFO: 既存のデータがないため、新規作成します。")
		state = &asset.State{
			UpdatedAt:       time.Now(),
			BaseCurrency:    "JPY",
			SelectedPhaseID: "phase_1",
			Phases:          make(map[string]asset.Phase),
			Assets:          make(map[string]asset.UserAsset),
		}

		state.Phases["phase_1"] = asset.Phase{
			Name:    "デフォルトフェーズ",
			Targets: make(map[string]float64),
		}
	}

	// 3. 引数がある場合は登録・更新
	if *owner != "" && *symbol != "" {
		_, exists := symbolMap[*symbol]

		if !exists && *symbol != "JPY" {
			log.Fatalf("❌ エラー: シンボル '%s' はスプレッドシートに存在しません。", *symbol)
		}

		if _, ok := state.Assets[*owner]; !ok {
			state.Assets[*owner] = make(asset.UserAsset)
		}

		currentAmount := state.Assets[*owner][*symbol]
		newAmount := currentAmount + *amount

		if newAmount <= 0 {
			delete(state.Assets[*owner], *symbol)
			fmt.Printf("🗑️ %s の %s を削除しました。\n", *owner, *symbol)
		} else {
			state.Assets[*owner][*symbol] = newAmount
			fmt.Printf("✅ %s に %s: %.4f を登録・更新しました。\n", *owner, *symbol, newAmount)
		}

		if len(state.Assets[*owner]) == 0 {
			delete(state.Assets, *owner)
		}

		state.UpdatedAt = time.Now()

		if err := io.SaveState(filePath, state); err != nil {
			log.Fatal("❌ データの保存に失敗しました:", err)
		}
	}

	// 4. 評価と計算
	rates := make(asset.Rates)
	for sym, info := range symbolMap {
		rates[sym] = info.Rate
	}

	evaluated := state.Evaluate(rates)

	// ================================================================================
	// 1. ヘッダー（選択中のフェーズ情報）の出力
	// ================================================================================
	targetPhase, phaseExists := state.Phases[state.SelectedPhaseID]
	if phaseExists {
		fmt.Println("================================================================================")
		fmt.Printf("%s%s\n", i18n[*lang]["target_phase"], targetPhase.Name)
		fmt.Printf("%s%s\n", i18n[*lang]["scenario"], targetPhase.Scenario)
		fmt.Printf("%s%s\n", i18n[*lang]["action"], targetPhase.Action)
		fmt.Println("================================================================================")
	}

	fmt.Println()
	fmt.Printf("%s: %.0f 円\n", i18n[*lang]["total_asset"], evaluated.Total)
	fmt.Println()

	// ================================================================================
	// 2. 資産配分と目標フェーズとの差分（テーブル出力）
	// ================================================================================
	fmt.Println(i18n[*lang]["diff_title"])
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(i18n[*lang]["table_header"])
	fmt.Println("--------------------------------------------------------------------------------")

	for symbol, diff := range evaluated.RebalanceGuide {
		statusStr := i18n[*lang]["status_opt"]
		if diff.DiffAmount > 1 {
			statusStr = i18n[*lang]["status_short"]
		} else if diff.DiffAmount < -1 {
			statusStr = i18n[*lang]["status_excess"]
		}

		diffAmountAbs := diff.DiffAmount
		if diffAmountAbs < 0 {
			diffAmountAbs = -diffAmountAbs
		}

		signStr := "+"
		if diff.DiffAmount < 0 {
			signStr = "-"
		} else if diff.DiffAmount == 0 {
			signStr = " "
		}

		fmt.Printf(" %-10s |  %5.1f %%  |  %5.1f %% | %s%.0f 円 (%s)\n",
			symbol,
			diff.CurrentPercent,
			diff.TargetPercent,
			signStr,
			diffAmountAbs,
			statusStr,
		)
	}
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println()

	// ================================================================================
	// 3. NEXT ACTIONS（誰が、どこで、何をすべきか）
	// ================================================================================
	fmt.Println(i18n[*lang]["next_actions"])
	fmt.Println("--------------------------------------------------------------------------------")

	actions := evaluated.GenerateNextActions(state)
	if len(actions) == 0 {
		fmt.Println(i18n[*lang]["no_rebalance"])
	} else {
		fmt.Println(i18n[*lang]["sell"])
		sellCount := 0
		for _, act := range actions {
			if act.Type == "SELL" {
				// 「推奨口座（任意）」というデータが入ってきたら、辞書から適切な言語に変換する
				owner := act.Owner
				if owner == "推奨口座（任意）" {
					owner = i18n[*lang]["default_owner"]
				}
				
				// 辞書からフォーマット（%s や %.0f の並び）を呼び出して出力
				fmt.Printf(i18n[*lang]["sell_action"], owner, act.Symbol, act.Amount)
				sellCount++
			}
		}
		if sellCount == 0 {
			fmt.Println(i18n[*lang]["none"])
		}

		fmt.Println()

		fmt.Println(i18n[*lang]["buy"])
		buyCount := 0
		for _, act := range actions {
			if act.Type == "BUY" {
				// 購入側も同様に、「推奨口座（任意）」を辞書から変換する
				owner := act.Owner
				if owner == "推奨口座（任意）" {
					owner = i18n[*lang]["default_owner"]
				}

				fmt.Printf(i18n[*lang]["buy_action"], owner, act.Symbol, act.Amount)
				buyCount++
			}
		}
		if buyCount == 0 {
			fmt.Println(i18n[*lang]["none"])
		}
	}
	fmt.Println("--------------------------------------------------------------------------------")
}