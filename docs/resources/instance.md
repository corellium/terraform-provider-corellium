# corellium_v1instance

## Example

```terraform
resource "corellium_v1instance" "example" {
  name = "example"
  flavor = "iphone7plus"
  project = "00000000-0000-4000-0000-000000000000"
  os = "15.7.5"
}
```

## Schema

### Required

- `name` (string) - The name of the instance.

- `flavor` (string) - The flavor of the instance.

- `project` (string) - The project ID of the instance.

- `os` (string) - The OS version of the instance.

### Optional

- `state` (string) - The state of the instance. Must be "on", "off" or "paused".

### Read-only

- `id` (string) - The ID of the instance.

- `key` (string) - The key of the instance.

- `state` (string) - The state of the instance. Possible to "on", "off", "paused", "creating", "deleting".

- `state_changed` (string) - The state change of the instance.

- `started_at` (string) - The start time of the instance.

- `user_task` (string) - The user task of the instance.

- `task_state` (string) - The task state of the instance.

- `error` (string) - The error of the instance.

- `boot_options` (object of `boot_options`) - The boot options of the instance.

- `service_ip` (string) - The service IP of the instance.

- `wifi_ip` (string) - The WiFi IP of the instance.

- `secondary_ip` (string) - The secondary IP of the instance.

- `services` (object of `services`) - The services of the instance.

- `panicked` (bool) - Whether the instance has panicked.

- `created` (bool) - Whether the instance has been created.

- `model` (string) - The model of the instance.

- `fwpackage` (string) - The firmware package of the instance.

- `os` (string) - The OS version of the instance.

- `agent` (object of `agent`) - The agent of the instance.

- `netnom` (object of `netnom`) - The netnom of the instance.

- `expose_port` (string) - The expose port of the instance.

- `fault` (string) - The fault of the instance.

- `patches` (list of string) - The patches of the instance.

- `created_by` (object of `user`) - The user who created the instance.

### Nested schema for `boot_options`

#### Read-only

- `boot_args` (string) - The boot args of the instance.

- `restore_boot_args` (string) - The restore boot args of the instance.

- `udid` (string) - The udid of the instance.

- `ecid` (string) - The ecid of the instance.

- `random_seed` (string) - The random seed of the instance.

- `pac` (bool) - Whether the instance has pac.

- `aprr` (bool) - Whether the instance has aprr.

- `additional_tags` (list of string) - The additional tags of the instance. Possible to "kalloc", "gpu", "no-keyboard", "nodevmode", "sep-cons-ext", "iboot-jailbreak", "llb-jailbreak", "rom-jailbreak".

### Nested schema for `services`

#### Read-only

- `vpn` (object of `vpn`) - The VPN of the instance.

### Nested schema for `vpn`

#### Read-only

- `proxy` (map) - The proxy of the instance.
- `listeners` (map) - The listeners of the instance.

### Nested schema for `agent`

#### Read-only

- `hash` (string) - The agent hash of the instance.

- `info` (string) - The agent info of the instance.

### Nested schema for `netnom`

#### Read-only

- `hash` (string) - The netmon hash of the instance.

- `info` (string) - The netmon info of the instance.

- `enabled` (bool) - Whether the instance has netnom enabled.

### Nested schema for `created_by`

### Read-only

- `id` (string) - The ID of the user.

- `username` (string) - The username of the user.

- `label` (string) - The label of the user.

- `deleted` (bool) - Whether the user has been deleted.
