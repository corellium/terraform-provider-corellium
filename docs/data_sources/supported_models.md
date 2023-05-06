# corellium_v1supportedmodels

## Example

```terraform
data "corellium_v1supportedmodels" "example" {}
```

## Schema

### Read-only

- `id` (string) - Supported models ID.

- `supported_models` (list of `supported_models`) - List of supported models.

### Nested schema for `supported_models`

#### Required

- `type` (string) - Model type.

- `name` (string) - Model name.

- `model` (string) - Model identifier.

- `flavor` (string) - Model flavor.

#### Optional

- `description` (string) - Model description.

- `board_config` (string) - Model board configuration.

- `platform` (string) - Model platform.

- `cp_id` (number) - Model CP ID.

- `bd_id` (number) - Model BD ID.

- `peripherals` (bool) - Whether the model has peripherals.