# corellium_v1team

## Example

```terraform
resource "corellium_v1team" "example" {
  label = "example"
  users = [
    {
      id = "00000000-0000-4000-0000-000000000000"
    },
  ]
}
```

## Schema

### Required

- `label` (string) - Team label.

### Optional

- `users` (list of `user`) - List of users to add to the team. Each user must have an `id` attribute.

### Read-only

- `id` (string) - Team ID.

### Nested schema for `user`

#### Read-only

- `id` (string) - User ID.
