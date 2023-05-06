# corellum_v1modelsoftware

## Example Usage

```terraform
data "corellium_v1modelsoftware" "example" {
  model = "iPhone15,3"
}
```

## Schema

### Required

- `model` (string) - Model identifier.

### Read-only

- `id` (string) - Model software ID.

- `model_software` (list of `model_software`) - List of model software.

### Nested schema for `model_software`

#### Optional

- `api_version` (string) - Android only API version.

- `android_flavor` (string) - Android only flavor.

- `build_id` (string) - Build ID.

- `filename` (string) - Filename.

- `md5_sum` (string) - MD5 sum.

- `orig_url` (string) - URL firmware is available at from vendor.

- `release_date` (string) - Release Date.

- `sha1_sum` (string) - SHA1 sum.

- `sha256_sum` (string) - SHA256 sum.

- `size` (number) - Size.

- `unique_id` (string) - Unique ID.

- `upload_date` (string) - Date uploaded.

- `url` (string) - URL firmware is available at.

- `version` (string) - Version.
