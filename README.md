# unleash-checker-ai

This tool identifies unused feature flags by referencing the Unleash API, extracts the relevant code sections, and uses LLM to correct the code.

It also targets potentially stale flags based on their lifetime and modifies the code accordingly.

ref: https://docs.getunleash.io/reference/technical-debt

Please note that using the OpenAI API for code modification will incur charges.

## Usage

Unleash Checker AI is intended to be used with GitHub Actions.

Please set the following environment variables in Actions Secret:

* UNLEASH_API_ENDPOINT: Unleash endpoint (https://app.unleash-hosted.com/api)
* UNLEASH_API_TOKEN: Unleash API token
* UNLEASH_PROJECT_ID: Project ID ("default")
* OPENAI_API_KEY: OpenAI API key

Set up a workflow like the one below in the repository you want to scan:


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
      - uses: kromiii/unleash-checker-ai@v0.1.12
        with:
          unleash_api_endpoint: ${{ secrets.UNLEASH_API_ENDPOINT }}
          unleash_api_token: ${{ secrets.UNLEASH_API_TOKEN }}
          unleash_project_id: ${{ secrets.UNLEASH_PROJECT_ID }}
          openai_api_key: ${{ secrets.OPENAI_API_KEY }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
          target_path: 'app'
```


You can narrow down the execution folder with `target_path`. By default, all files are targeted, but it is recommended to specify it to exclude third-party libraries, etc.

Here is an example of the generated PR

https://github.com/kromiii/unleash-checker-ai/pull/13

If you are using GHES, please add `GITHUB_BASE_URL` to the Actions parameters.

```yaml
        with:
          unleash_api_endpoint: ${{ secrets.UNLEASH_API_ENDPOINT }}
          unleash_api_token: ${{ secrets.UNLEASH_API_TOKEN }}
          unleash_project_id: ${{ secrets.UNLEASH_PROJECT_ID }}
          openai_api_key: ${{ secrets.OPENAI_API_KEY }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
          target_path: 'app'
          github_base_url: 'https://git.example.com'
```
