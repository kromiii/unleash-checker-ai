# unleash-checker-ai

unleash api を参照して使われていないflagを発見、コードから該当する部分を抽出し、LLMによるコードの修正を行う script です

コードの修正に openai api を使用しているため課金が発生します

## 使い方

環境変数に以下の情報を設定

```
export UNLEASH_API_ENDPOINT=http://localhost:4242/api
export UNLEASH_API_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxx
export UNLEASH_PROJECT_ID=default
export OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxx

```

package をビルド

```
go build ./cmd/unleash-checker-ai
```

対象のフォルダを指定して実行

```
./unleash-checker-ai example
```

## 実行結果

実行すると以下のような結果が出力されます

```
% ./unleash-checker-ai example     
Stale or potentially stale flags:
 - unleash-ai-example-stale
These flags are used in:
 - example/hello_world_stale.py
Removing unused flags by LLM...
Done!
```

gitで差分を見つつ、適宜PRなど立ててください

## ディレクトリ構造

```
unleash-checker-ai/
├── cmd/
│   └── unleash-checker-ai/
│       └── main.go
├── internal/
│   ├── unleash/
│   │   └── client.go
│   ├── finder/
│   │   └── finder.go
│   ├── modifier/
│   │   └── modifier.go
│   └── config/
│       └── config.go
├── pkg/
│   └── openai/
│       └── client.go
├── go.mod
├── go.sum
└── README.md
```
