# corellium_v1roles

## Example

```terraform
data "corellium_v1roles" "example" {
  project = "00000000-0000-4000-0000-000000000000"
}
```

## Schema

### Optional

- `project` (string) - The project ID to list roles for.

### Read-only

- `roles` (list of `role`) - The list of role.

### Nested schema for `role`

### Optional

- `project` (string) - The project ID.

- `user` (string) - The user ID.

### Read-only

- `role` (string) - The role name. Possible to be: "admin" or "_member_".