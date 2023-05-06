# corelium_v1getinstance

## Example

```terraform
data "corellium_v1getinstance" "example" {}
```

## Schema

### Optional

- `instances` (list of `instance` objects)

### Nested schema for `instance`

#### Optional

- `id` (string) - Instance ID.

- `name` (string) - Instance name.

- `key` (string) - Instance key.

- `flavor` (string) - Instance flavor.

- `type` (string) - Instance type.

- `project` (string) - Instance project.

- `state` (string) - Instance state.

- `state_changed` (string) - Instance state changed.

- `started_at` (string) - Instance started at.

- `user_task` (string) - Instance user task.

- `task_state` (string) - Instance task state.

- `error` (string) - Instance error.

- `service_ip` (string) - Instance service IP.

- `wifi_ip` (string) - Instance wifi IP.

- `secondary_ip` (string) - Instance secondary IP.

- `panicked` (bool) - Instance panicked.

- `created` (string) - Instance created.

- `model` (string) - Instance model.

- `fwpackage` (string) - Instance fwpackage.

- `os` (string) - Instance os.

- `netmon` (object of `netmon`) - Instance netmon.

- `expose_port` (string) - Instance expose port.

- `fault` (bool) - Instance fault.

- `patches` (list of strings) - Instance patches.

- `created_by` (object of `created_by`) - User who created the instance.

- `boot_options` (object of `boot_optons`) - Boot options for the instance.

- `agent` (object of `agent`) - Agent information for the instance.

### Nested schema for `netmon`

#### Optional

- `hash` (string) - Netmon hash.

- `info` (string) - Netmon info.

- `enabled` (bool) - Netmon enabled.

### Nested schema for `created_by`

#### Optional

- `id` (string) - User ID.

- `username` (string) - User username.

- `label` (string) - User label.

- `deleted` (bool) - User deleted.

### Nested schema for `boot_options`

#### Optional

- `boot_args` (string) - Boot args.

- `restore_boot_args` (string) - Restore boot args.

- `udid` (string) - UDID.

- `ecid` (string) - ECID.

- `random_seed` (string) - Random seed.

- `pac` (bool) - PAC.

- `aprr` (bool) - APRR.

- `additional_tags` (list of strings) - Additional tags.

### Nested schema for `agent`

#### Optional

- `hash` (string) - Agent hash.

- `info` (string) - Agent info.
