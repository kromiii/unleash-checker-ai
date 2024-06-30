# unleash-checker-ai

unleash api を参照して使われていないflagを発見、コードから該当する部分を抽出し、LLMによるコードの修正を行う script です

## 使い方

環境変数に以下の情報を設定

```
export UNLEASH_API_ENDPOINT=
export UNLEASH_API_TOKEN=
export OPENAI_API_KEY=
```

build

```
go build ./cmd/unleash-checker-ai
```

対象とするフォルダを指定

```
unleash-checker-ai ~/src/example.com/repo
```

使われていない flag をリストアップし

```
unused flags:
 - flag_1
 - flag_2
```

その flag が使われているファイルを列挙

```
These flags are used in:
 - hoge.py
 - fuga.py
```

それぞれのファイルについて unleash flag を取り除いたファイルで上書きします（この過程で openai api へ接続します）

```
unused flags are removed 
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
