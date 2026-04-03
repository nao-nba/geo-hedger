package asset

import (
	"time"
)

// スプレッドシートから取得した最新レートを保持する型
// 例: {"USD": 150.5, "GLD": 30000}
type Rates map[string]float64

type SymbolInfo struct {
	Rate     float64
	Category string
}

// ★新規追加: シンボル名をキーとした詳細情報のマップ
type SymbolMap map[string]SymbolInfo

// 正常にレートが引けるかチェックする補助関数
func (r Rates) GetRate(symbol string) (float64, bool) {
	rate, ok := r[symbol]
	return rate, ok
}

// フェーズの定義：自分の想定するシナリオをフェーズに区切り、その理由をログとして残す
type Philosophy struct {
	CurrentPhase string `json:"current_phase,omitempty"` // 例: "Phase1_MiddleInflation"
	Rationale    string `json:"rationale,omitempty"`     // 理由の言語化。ここに思考のログを残す
}

// 各フェーズにおける理想の資産比率（%）を定義する
type TargetAllocation struct {
	Currencies  float64 `json:"currencies,omitempty"` // 例：20.0(%)
	WesternRisk  float64 `json:"western_risk,omitempty"` // 例：30.0(%)
	EasternRisk  float64 `json:"eastern_risk,omitempty"` // 例：20.0(%)
	StatelessRisk  float64 `json:"stateless_risk,omitempty"` // 例：30.0(%)
}

// 現在の実際の資産の保有数量（口数や枚数）
type Portfolio struct {
	Currencies  map[string]float64 `json:"currencies"` // "JPY": 10000, "USD": 5000
	WesternRisk map[string]float64 `json:"western_risk"` // "GLD": 4000
	EasternRisk map[string]float64 `json:"eastern_risk"` // "hsi_2800": 1000
	StatelessRisk map[string]float64 `json:"stateless_risk"` // "BTC": 1.5
}

// 現在の実際の資産のレート換算後の「円建ての金額」。レートはGoogleスプレッドシートから取得
type EvaluatedPortfolio struct {
	Currencies    map[string]float64 `json:"currencies"`     // "USD": 150000 (円)
	WesternRisk   map[string]float64 `json:"western_risk"`   // "GLD": 300000 (円)
	EasternRisk   map[string]float64 `json:"eastern_risk"`   // "1321": 200000 (円)
	StatelessRisk map[string]float64 `json:"stateless_risk"` // "BTC": 4500000 (円)
}

// State は資産管理の全体的なデータ構造（静的なデータ）
type State struct {
	UpdatedAt time.Time  `json:"updated_at"` // 最終更新日時
	Mindset   Philosophy `json:"mindset"`    // 現在の地政学フェーズと理由
	Targets TargetAllocation `json:"targets"` // 理想の比率
	Assets    Portfolio  `json:"assets"`     // 実際の資産データ（保有数量）
}

// TODO: 比率の計算や差分を出す「関数（ロジック）」を後から追加していく

// EvaluatedPortfolio（円建て評価額）を計算するロジック（Portfolio*Googleスプレッドシートのレート）
func (p *Portfolio) Evaluate(rates map[string]float64) *EvaluatedPortfolio {
	eval := &EvaluatedPortfolio{
		Currencies:    make(map[string]float64),
		WesternRisk:   make(map[string]float64),
		EasternRisk:   make(map[string]float64),
		StatelessRisk: make(map[string]float64),
	}

	// 通貨の換算
	for code, amount := range p.Currencies {
		if rate, ok := rates[code]; ok {
			eval.Currencies[code] = amount * rate
		} else if code == "JPY" {
			eval.Currencies[code] = amount // JPYはそのまま
		}
	}

	// 西側リスクの換算
	for code, amount := range p.WesternRisk {
		if rate, ok := rates[code]; ok {
			eval.WesternRisk[code] = amount * rate
		}
	}

	// 東側リスクの換算
	for code, amount := range p.EasternRisk {
		if rate, ok := rates[code]; ok {
			eval.EasternRisk[code] = amount * rate
		}
	}

	// 無国籍リスクの換算
	for code, amount := range p.StatelessRisk {
		if rate, ok := rates[code]; ok {
			eval.StatelessRisk[code] = amount * rate
		}
	}

	return eval
}

// Portfolioの中にあるすべての資産（円換算）の合計を計算
// レート換算後の値をEvaluatedPortfolioと仮定。換算する工程を後日実装する
func (p *EvaluatedPortfolio) CalculateTotal() float64 {
	total := 0.0

	// 通貨（Currencies）の合計
	for _, amount := range p.Currencies {
		total += amount
	}
	// 西側リスク（WesternRisk）の合計
	for _, amount := range p.WesternRisk {
		total += amount
	}
	// 東側リスク（EasternRisk）の合計
	for _, amount := range p.EasternRisk {
		total += amount
	}
	// 無国籍リスク（StatelessRisk）の合計
	for _, amount := range p.StatelessRisk {
		total += amount
	}

	return total
}