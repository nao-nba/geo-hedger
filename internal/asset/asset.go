package asset

import (
	"time"
)

// スプレッドシートから取得した最新レートを保持する型
type Rates map[string]float64

type SymbolInfo struct {
	Rate     float64
	Category string
}

// シンボル名をキーとした詳細情報のマップ
type SymbolMap map[string]SymbolInfo

// 1. モデル定義の刷新 --------------------------------------------------

// State は資産管理の全体的なデータ構造（state.json に対応）
type State struct {
	UpdatedAt       time.Time            `json:"updated_at"`
	BaseCurrency    string               `json:"base_currency"`      // 追加: 基準通貨（例: "JPY"）
	SelectedPhaseID string               `json:"selected_phase_id"`  // 追加: 比較したいフェーズのID
	Phases          map[string]Phase     `json:"phases"`             // 追加: 複数のフェーズを保持
	Assets          map[string]UserAsset `json:"assets"`             // 変更: 家族単位での資産保有数量
}

// Phase は地政学フェーズの定義
type Phase struct {
	Name     string             `json:"name"`
	Scenario string             `json:"scenario"`
	Action   string             `json:"action"`
	Risk     string             `json:"risk,omitempty"` // 任意項目
	Targets  map[string]float64 `json:"targets"`        // 理想の比率（例: {"JPY": 30.0, "BTC": 30.0}）
}

// UserAsset は特定の持ち主（家族）が持つ資産のマップ
// 例: {"JPY": 1000000, "BTC": 0.2}
type UserAsset map[string]float64

// 現在の比率と、目標との乖離
type AllocationDiff struct {
	CurrentPercent float64 `json:"current_percent"` // 現在の比率 (%)
	TargetPercent  float64 `json:"target_percent"`  // 目標の比率 (%)
	DiffAmount     float64 `json:"diff_amount"`     // 金額ベースの過不足（基準通貨建て）
}

// 評価結果の全体像
type EvaluationResult struct {
	Total          float64                   `json:"total"`            // 総資産
	Breakdown      map[string]float64        `json:"breakdown"`        // シンボル別の合計額
	RebalanceGuide map[string]AllocationDiff `json:"rebalance_guide"` // リバランスのガイド
}

// 2. 計算ロジック ----------------------------------------------------

// 実際の保有数量（口数）にスプレッドシートのレートを掛け合わせて、基準通貨建ての金額に変換する
func (s *State) Evaluate(rates Rates) *EvaluationResult {
	result := &EvaluationResult{
		Breakdown:      make(map[string]float64),
		RebalanceGuide: make(map[string]AllocationDiff),
	}

	total := 0.0

	// 家族全員の資産を合算し、基準通貨建てに換算する
	for _, userAsset := range s.Assets {
		for symbol, amount := range userAsset {
			evaluatedAmount := amount

			// 基準通貨（JPYなど）以外の場合は、レートを掛けて換算する
			if symbol != s.BaseCurrency {
				if rate, ok := rates[symbol]; ok {
					evaluatedAmount = amount * rate
				}
			}

			result.Breakdown[symbol] += evaluatedAmount
			total += evaluatedAmount
		}
	}

	result.Total = total

	// 選択されている目標フェーズを取得
	targetPhase, ok := s.Phases[s.SelectedPhaseID]
	if !ok {
		return result // フェーズが見つからない場合はここで終了
	}

	// 現在の比率の計算と、目標比率との乖離（ガイド）の作成
	for symbol, targetPercent := range targetPhase.Targets {
		currentAmount := result.Breakdown[symbol]
		currentPercent := 0.0
		if total > 0 {
			currentPercent = (currentAmount / total) * 100
		}

		targetAmount := total * (targetPercent / 100)
		diffAmount := targetAmount - currentAmount

		result.RebalanceGuide[symbol] = AllocationDiff{
			CurrentPercent: currentPercent,
			TargetPercent:  targetPercent,
			DiffAmount:     diffAmount,
		}
	}

	return result
}

// 3. アクションロジックの追加 --------------------------------------------

// DeleteAsset は特定の持ち主の特定の資産を削除、または一部削減します。
// amount に 0 を指定、または現在の保有量以上の数字を指定した場合は、その資産のデータを完全に削除します。
func (s *State) DeleteAsset(owner string, symbol string, amount float64) {
	// 指定された owner（家族）がいるかチェック
	userAsset, ok := s.Assets[owner]
	if !ok {
		return // 家族がいなければ何もしない
	}

	// 指定された資産を持っているかチェック
	currentAmount, ok := userAsset[symbol]
	if !ok {
		return // 資産を持っていなければ何もしない
	}

	// 数量の指定がない（0）、または現在の量より多く引こうとした場合は「完全削除」
	if amount == 0 || amount >= currentAmount {
		delete(userAsset, symbol)
	} else {
		// それ以外（一部削除）は引き算
		userAsset[symbol] = currentAmount - amount
	}

	// もしその家族の資産がゼロ（空マップ）になったら、家族のデータごと削除して整理
	if len(userAsset) == 0 {
		delete(s.Assets, owner)
	}
}

// NextAction は、リバランスのために具体的に「誰が、何を、いくら売買すべきか」の指示を生成します
type NextAction struct {
	Owner  string  // 誰の（例: "家族A（SBI証券）"）
	Symbol string  // 何を（例: "BTC"）
	Type   string  // "BUY"（購入） または "SELL"（売却）
	Amount float64 // 基準通貨建ての金額（例: 500000 円分）
}

// GenerateNextActions は乖離額から具体的な売買行動をリスト化します
func (r *EvaluationResult) GenerateNextActions(s *State) []NextAction {
	var actions []NextAction

	for symbol, diff := range r.RebalanceGuide {
		// 乖離額がゼロ、または微小（1円未満など）なら無視
		if diff.DiffAmount > -1 && diff.DiffAmount < 1 {
			continue
		}

		// 超過している場合（売却指示）
		if diff.DiffAmount < 0 {
			excessAmount := -diff.DiffAmount // プラスの値に変換
			
			// 家族の資産を走査して、そのシンボルを持っている人から引いていく
			for owner, userAsset := range s.Assets {
				if _, ok := userAsset[symbol]; ok {
					// 実際の保有額を概算（簡略化のため現在のレートで換算）
					// 本来は厳密なレート計算が必要ですが、まずはロジックを回すために金額ベースで割り振ります
					actions = append(actions, NextAction{
						Owner:  owner,
						Symbol: symbol,
						Type:   "SELL",
						Amount: excessAmount, // 一旦、過不足額をそのまま割り当て（TODO: 家族の保有量上限でキャップするロジック）
					})
					break // ひとまず1人見つけたら割り当てて次へ（簡易版）
				}
			}
		}

		// 不足している場合（購入指示）
		if diff.DiffAmount > 0 {
			// 購入は、そのシンボルを既に持っている、あるいは代表して家族Aに割り振る
			// ここではシンプルに「目標フェーズでそのシンボルを割り当てられている代表的な口座」に紐付けるか、
			// ひとまずアクションのリストとして「このシンボルがこれだけ不足している」と出す形にします
			actions = append(actions, NextAction{
				Owner:  "推奨口座（任意）",
				Symbol: symbol,
				Type:   "BUY",
				Amount: diff.DiffAmount,
			})
		}
	}

	return actions
}