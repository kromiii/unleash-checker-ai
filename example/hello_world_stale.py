import unleash_client

client = unleash_client.initialize_client(
    url="https://your-unleash-instance.com/api/",
    client_keys={"default": "your-client-side-api-key"},
    app_name="your-app-name"
)

if client.is_enabled("unleash-ai-example-stale"):
    print("hello world")

client.close()
