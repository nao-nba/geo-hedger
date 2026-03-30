package asset

import (
	"time"
)

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

// 現在の実際の資産の額（円換算など）を定義
type Portfolio struct {
	Currencies  map[string]float64 `json:"currencies"` // "JPY": 10000, "USD": 5000
	WesternRisk map[string]float64 `json:"western_risk"` // "GLD": 4000
	EasternRisk map[string]float64 `json:"eastern_risk"` // "hsi_2800": 1000
	StatelessRisk map[string]float64 `json:"stateless_risk"` // "BTC": 1.5
}

// State は資産管理の全体的なデータ構造（大枠）
type State struct {
	UpdatedAt time.Time  `json:"updated_at"` // 最終更新日時
	Mindset   Philosophy `json:"mindset"`    // 現在の地政学フェーズと理由
	Targets TargetAllocation `json:"targets"` // 理想の比率
	Assets    Portfolio  `json:"assets"`     // 実際の資産データ
}

// TODO: 比率の計算や差分を出す「関数（ロジック）」を後から追加していく