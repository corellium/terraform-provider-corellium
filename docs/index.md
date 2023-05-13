# Corellium provider

## Example

```terraform
provider "corellium" {
  token = ""
  host = "app.corellium.com"
}
```

## Schema

### Required

- `token` (string) - Corellium API token. This can also be set via the CORELLIUM_API_TOKEN environment variable.

### Optional

- `host` (string) - Corellium API host. This can also be set via the CORELLIUM_API_HOST environment variable. Default value is `app.corellium.com".
