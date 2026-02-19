# WODATA Plan

## Goal

Create a Terraform provider that functions exactly like the built-in `terraform_data` resource, with one addition: support for write-only input arguments (`input_wo` and `input_wo_version`).

This enables capturing ephemeral values (like public keys from ephemeral `tls_private_key`) into persistent state without storing the corresponding private key in state.

## Provider Name

**wodata**

## Resource Name

**wo_data**

## API

```hcl
resource "wodata_wo_data" "example" {
  # Write-only input - accepts ephemeral values
  input_wo         = ephemeral.tls_private_key.ssh.public_key_openssh
  input_wo_version = 1

  # Optional: triggers_replace (same as terraform_data)
  triggers_replace = [some_resource.id]
}

# Output - readable from state
output "public_key" {
  value = wodata_wo_data.example.output
}
```

## Arguments

- `input_wo` (String, Optional, Write-only) - Accepts ephemeral values; value is stored in `output`
- `input_wo_version` (Number, Optional) - When changed, triggers re-capture of `input_wo` value
- `triggers_replace` (Any, Optional) - Forces replacement when value changes

## Attributes

- `id` (String, Computed) - Unique identifier
- `output` (Any, Computed) - The captured value from `input_wo`
