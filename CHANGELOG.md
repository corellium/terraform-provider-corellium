# Changelog

## 0.1.0-alpha (Unreleased)

This release switches the provider to the public Corellium API client
([corellium-go-api-client](https://github.com/corellium/corellium-go-api-client))
and publishes it under the `corellium/corellium` registry address. It contains
breaking changes, and state from 0.0.x releases can't be upgraded in place; see
[Upgrading from 0.0.x](#upgrading-from-00x).

### Breaking changes

- **Provider address:** the provider is now `corellium/corellium` (was
  `aimoda/corellium`). Update `required_providers`:

  ```hcl
  terraform {
    required_providers {
      corellium = {
        source = "corellium/corellium"
      }
    }
  }
  ```

- **Terraform version:** the provider now serves plugin protocol 6, which
  requires Terraform 1.0 or later.

- **Corellium API version:** the provider now calls the Corellium v3 API
  (`/api/v3/...`). Servers that only expose the v1 API are not supported.

- **Resources renamed:** the `v1` prefix was dropped from all resource types.

  | 0.0.x                   | 0.1.0                 |
  |-------------------------|-----------------------|
  | `corellium_v1image`     | `corellium_image`     |
  | `corellium_v1instance`  | `corellium_instance`  |
  | `corellium_v1project`   | `corellium_project`   |
  | `corellium_v1snapshot`  | `corellium_snapshot`  |
  | `corellium_v1team`      | `corellium_team`      |
  | `corellium_v1user`      | `corellium_user`      |
  | `corellium_v1webplayer` | `corellium_webplayer` |

- **Data sources renamed:** the `v1` prefix was also dropped from all data
  source types.

  | 0.0.x                         | 0.1.0                       |
  |-------------------------------|-----------------------------|
  | `corellium_v1instances`       | `corellium_instances`       |
  | `corellium_v1modelsoftware`   | `corellium_modelsoftware`   |
  | `corellium_v1projects`        | `corellium_projects`        |
  | `corellium_v1ready`           | `corellium_ready`           |
  | `corellium_v1roles`           | `corellium_roles`           |
  | `corellium_v1supportedmodels` | `corellium_supportedmodels` |

  Data sources hold no state, so renaming them in the configuration is enough.

- **`corellium_instance`:** the `os` attribute was replaced by two required
  attributes, `os_version` and `os_build`:

  ```hcl
  # 0.0.x
  os = "16.0.0"

  # 0.1.0
  os_version = "16.0.0"
  os_build   = "20A362"
  ```

  The same change applies to each instance in the `corellium_instances` data
  source.

- **`corellium_image`:** the computed `unique_id` attribute was removed.

### Upgrading from 0.0.x

There is no automatic state upgrade: the resource types and provider address
both changed, so Terraform treats 0.1.0 resources as new, and the resources
don't support `terraform import` yet. To move an existing configuration:

1. With the 0.0.x provider still installed, run `terraform destroy`, or
   `terraform state rm` the resources you want to keep and manage them outside
   Terraform.
2. Update the configuration: the provider `source`, resource and data source
   type names, and the instance `os` attribute, as described above.
3. Run `terraform init -upgrade` and `terraform apply` to create the resources
   under 0.1.0.
