# Writeonly Provider Plan

## Goal

Create a Terraform provider that functions exactly like the built-in `terraform_data` resource, with one addition: support for write-only input arguments (`input_wo` and `input_wo_version`).

This enables capturing ephemeral values (like public keys from ephemeral `tls_private_key`) into persistent state without storing the corresponding private key in state.

## Provider Name

**writeonly**

## Resource Name

**writeonly_data**

## API

```hcl
provider "writeonly" {}

resource "writeonly_data" "example" {
  # Write-only input - accepts ephemeral values
  input_wo         = ephemeral.tls_private_key.ssh.public_key_openssh
  input_wo_version = 1

  # Optional: triggers_replace (same as terraform_data)
  triggers_replace = [some_resource.id]
}

# Output - readable from state
output "public_key" {
  value = writeonly_data.example.output
}
```

## Arguments

- `input_wo` (String, Optional, Write-only) - Accepts ephemeral values; value is stored in `output`
- `input_wo_version` (Number, Optional) - When changed, triggers re-capture of `input_wo` value
- `triggers_replace` (Any, Optional) - Forces replacement when value changes

## Attributes

- `id` (String, Computed) - Unique identifier
- `output` (Any, Computed) - The captured value from `input_wo`

## Not Supported

- **Import**: Not supported, matching the behavior of `terraform_data` (the built-in resource this provider mimics). There is no external infrastructure to import - the `output` value is derived from config.

## Write Only Values

This provider demonstrates the write-only arguments pattern introduced in Terraform 1.11.

### Schema Implementation

```go
// Schema definition for write-only attribute
"input_wo": schema.StringAttribute{
    Optional: true,
    WriteOnly: true,  // This makes it write-only
}

// Version argument (required to track changes)
"input_wo_version": schema.Int64Attribute{
    Optional: true,
    Computed: true,
    PlanModifiers: []planmodifier.Int64{
        int64planmodifier.UseStateForUnknown(),
    },
}
```

### Key Implementation Details

1. **Access from config, not plan**: Write-only values only exist in `req.Config`, never in `req.Plan` or `req.State`
   ```go
   diags := req.Config.Get(ctx, &data)  // Correct
   // NOT req.Plan.Get(ctx, &data) - write-only values are null in plan
   ```

2. **Never persisted**: The framework automatically nullifies write-only values before sending to state/plan. The provider must capture the value to a separate computed attribute (`output`) if it needs to be stored.

3. **Version triggers updates**: Since Terraform cannot track write-only value changes, incrementing `input_wo_version` triggers the provider to re-read `input_wo` from config.

4. **Accepts ephemeral values**: Write-only arguments can receive values from `ephemeral` resources, which would otherwise be impossible to persist.

### Usage with Ephemeral Resources

```hcl
# Generate ephemeral key pair
ephemeral "tls_private_key" "example" {
  algorithm = "ED25519"
}

# Capture the public key into persistent state
resource "writeonly_data" "example" {
  input_wo         = ephemeral.tls_private_key.example.public_key_openssh
  input_wo_version = 1
}

# Now readable from state
output "public_key" {
  value = writeonly_data.example.output
}
```

### For Future Developers

- Use `terraform-plugin-framework` (not SDKv2) for write-only support
- Always pair write-only arguments with a version argument
- Store captured values in computed attributes to persist in state
- Test with ephemeral resources to verify the pattern works correctly
