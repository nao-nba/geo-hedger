package io

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nao-nba/geo-hedger/internal/asset"
)

// LoadState は、指定されたパスのJSONファイルを読み込んで Goのデータ（Stateモデル）に変換
func LoadState(filePath string) (*asset.State, error) {
	// 1. ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 2. JSONをGoの構造体（State）に変換（デコード）する
	var state asset.State
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
}

// SaveState は、Goのデータ（asset.State 構造体のデータ）をJSONファイルに変換
func SaveState(filePath string, state *asset.State) error {
	// 1. 保存先のディレクトリが存在しない場合は作成する
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 2. ファイルを書き込みモードで開く（なければ作成）
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3. 構造体をJSONに変換（エンコード）してファイルに書き込む
	encoder := json.NewEncoder(file)
	// インデント（字下げ）をつけて、人間がエディタで読みやすい綺麗なJSONにする
	encoder.SetIndent("", "  ") 
	
	if err := encoder.Encode(state); err != nil {
		return err
	}

	return nil
}