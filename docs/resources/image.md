# corellium_v1image

## Example

```terraform
resource "corellium_v1image" "example" {
  name = "example"
  type = "backup"
  filename = "/tmp/image.txt"
  encapsulated = false
  project = "00000000-0000-4000-0000-000000000000"
}
```

## Schema

### Required

- `name` (string) - Image name.

- `type` (string) - Image type. Must be one of `fwbinary`, `kernel`, `devicetree`, `ramdisk`, `loaderfile`, `sepfw`, `seprom`, `bootrom`, `llb`, `iboot`, `ibootdata`, `fwpackage`, `partition`, or `backup`.

- `filename` (string) - Path to the image file.

- `encapsulated` (bool) - Whether the image is encapsulated.

- `project` (string) - Project ID.

### Read-only

- `status` (string) - Image status.

- `id` (string) - Image ID.

- `unique_id` (string) - Image unique ID.

- `size` (number) - Image size.

- `created_at` (string) - Image creation time.
