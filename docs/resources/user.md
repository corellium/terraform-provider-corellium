# corellium_v1user

## Example

```terraform
resource "corellium_v1user" "example" {
  name = "example"
  label = "example"
  email = "example@email.com"
  password = "examplepassword"
  administrator = false
}
```

## Schema

### Required

- `name` (string) - User name.

- `label` (string) - User label.

- `email` (string) - User email.

- `password` (string, sensitive) - User password.

- `administrator` (bool) - User administrator status.

### Read-only

- `id` (string) - User ID.
