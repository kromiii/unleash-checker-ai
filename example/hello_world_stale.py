import unleash_client
import unleash_client

client = unleash_client.initialize_client(
    url="https://your-unleash-instance.com/api/",
    client_keys={"default": "your-client-side-api-key"},
    app_name="your-app-name"
)

# The feature flag "unleash-ai-example-stale" was always enabled
print("hello world")

client.close()
