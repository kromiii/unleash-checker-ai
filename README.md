# unleash-checker-ai

unleash api を参照して使われていないflagを発見、コードから該当する部分を抽出し、LLMによるコードの修正を行う script です

## 使い方

環境変数に以下の情報を設定

```
export UNLEASH_API_ENDPOINT=
export UNLEASH_API_KEY=
```

対象とするフォルダを指定

```
unleash-checker-ai ~/src/example.com/repo
```

使われていない flag をリストアップし

その flag が使われているファイルを列挙

それぞれのファイルについて unleash flag を取り除いたファイルで上書きします

gitで差分を見つつ、適宜PRなど立ててください
