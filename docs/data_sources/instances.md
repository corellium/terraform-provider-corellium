# corelium_v1instances

## Example

```terraform
data "corellium_v1instances" "example" {}
```

## Schema

### Optional

- `id` (string) - ID of instances' list.

- `instances` (list of `instance` objects) list of instances.

### Nested schema for `instance`

#### Optional

- `id` (string) - The ID of the instance.

- `name` (string) - The name of the instance.

- `key` (string) - The key of the instance.

- `flavor` (string) - The flavor of the instance.
	A flavor is a device model, what can be a Android or iOS device.

	The following flavors are examples of supported flavors for Android:
    - ranchu (for Generic Android devices)
    - google-nexus-4
    - google-nexus-5
    - google-nexus-5x
    - google-nexus-6
    - google-nexus-6p
    - google-nexus-9
    - google-pixel
    - google-pixel-2
    - google-pixel-3
    - htc-one-m8
    - huawei-p8
    - samsung-galaxy-s-duos

	The following flavors are examples for iOS:
    - iphone6
    - iphone6plus
    - ipodtouch6
    - ipadmini4wifi
    - iphone6s
    - iphone6splus
    - iphonese
    - iphone7
    - iphone7plus
    - iphone8
    - iphone8plus
    - iphonex
    - iphonexs
    - iphonexsmax
    - iphonexsmaxww
    - iphonexr
    - iphone11
    - iphone11pro
    - iphone11promax
    - iphonese2
    - iphone12m
    - iphone12
    - iphone12p
    - iphone12pm
    - iphone13
    - iphone13m
    - iphone13p
    - iphone13pm

- `type` (string) - The type of the instnace.

- `project` (string) - The project ID of the instance.

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

- `panicked` (bool) - Whether the instance has panicked.

- `created` (bool) - Whether the instance has been created.

- `model` (string) - The model of the instance.

- `fwpackage` (string) - The firmware package of the instance.

- `os` (string) - The OS version of the instance.

- `agent` (object of `agent`) - The agent of the instance.

- `netnom` (object of `netnom`) - The netnom of the instance.

- `expose_port` (string) - The expose port of the instance.

- `fault` (string) - The fault of the instance.

- `patches` (list of string) - The patches of the instance. Possible to be "jailbroken", "nonjailbroken" or "corelliumd". "jailbroken" is the default value. "nonjailbroken" means that instance should not be jailbroken and "corelliumd", the instance should not be jailbroken, but should profile API agent.

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
