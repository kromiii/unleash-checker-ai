# unleash-checker-ai

unleash api を参照して使われていないflagを発見、コードから該当する部分を抽出し、LLMによるコードの修正を行うツールです

flag の lifetime から potentially stale な flag も対象とし、コードの修正を行います

ref: https://docs.getunleash.io/reference/technical-debt

コードの修正に openai api を使用しているため課金が発生します

トークンの長いファイルに対しては、LLMによる修正は行わず、コメントで flag の使用箇所を示すのみとなります

## 使い方

Unleash Checker AI は GitHub Actions での利用を想定しています

Actions Secret に以下の環境変数を設定してください

```
UNLEASH_API_ENDPOINT: Unleash のエンドポイント (https://app.unleash-hosted.com/api)
UNLEASH_API_TOKEN: Unleash の API トークン
UNLEASH_PROJECT_ID: プロジェクトID ("default")
OPENAI_API_KEY: OpenAI API キー
```

スキャンしたいレポジトリで以下のようなワークフローを設定してください

```yaml
name: Unleash Checker
on:
  schedule:
    - cron: '0 0 * * *'    
jobs:
  unleash_checker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: kromiii/unleash-checker-ai@v2
        with:
          unleash_api_endpoint: ${{ secrets.UNLEASH_API_ENDPOINT }}
          unleash_api_token: ${{ secrets.UNLEASH_API_TOKEN }}
          unleash_project_id: ${{ secrets.UNLEASH_PROJECT_ID }}
          openai_api_key: ${{ secrets.OPENAI_API_KEY }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
          target_path: 'app'
```

`target_path` でフォルダを絞って実行することができます。デフォルトは全てのファイルが対象となりますが、サードパーティのライブラリなどを除外するために指定することをお勧めします。

生成されるPRのサンプルはこちら

https://github.com/kromiii/sample-app-for-unleash-checker/pull/2
