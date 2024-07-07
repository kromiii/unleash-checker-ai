# unleash-checker-ai

unleash api を参照して使われていないflagを発見、コードから該当する部分を抽出し、LLMによるコードの修正を行う script です

flag の lifetime から potentially stale な flag も対象とし、コードの修正を行います

ref: https://docs.getunleash.io/reference/technical-debt

コードの修正に openai api を使用しているため課金が発生します

## 使い方

Unleash Checker AI は GitHub Actions での利用を想定しています

Actions Secret に以下の環境変数を設定してください

```
export UNLEASH_API_ENDPOINT=http://example.com/api
export UNLEASH_API_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxx
export UNLEASH_PROJECT_ID=default
export OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxx
```

以下のような workflow を設定してください

```yaml
name: Unleash Checker
on:
  workflow_dispatch:
    
jobs:
  unleash_checker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: kromiii/unleash-checker-ai@v0.1.6
        with:
          unleash_api_endpoint: ${{ secrets.UNLEASH_API_ENDPOINT }}
          unleash_api_token: ${{ secrets.UNLEASH_API_TOKEN }}
          unleash_project_id: ${{ secrets.UNLEASH_PROJECT_ID }}
          openai_api_key: ${{ secrets.OPENAI_API_KEY }}
          target_path: 'app'
```

`target_path` には対象のフォルダを指定してください

## ローカルでの実行

package をビルド

```
go build ./cmd/unleash-checker-ai
```

対象のフォルダを指定して実行

```
./unleash-checker-ai example
```

実行後はgitで差分を見つつ、適宜PRなど立ててください

## Option

`--only-stale` オプションを指定すると、potentially stale flags は無視されます
