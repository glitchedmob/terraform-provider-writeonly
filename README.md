# Terraform Provider Wodata

A Terraform provider that functions like the built-in `terraform_data` resource, with support for write-only input arguments. This enables capturing ephemeral values (like public keys from `tls_private_key`) into persistent state without storing the corresponding private key.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.10
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

```shell
go install
```

## Using the Provider

```hcl
resource "wodata_wo_data" "example" {
  # Write-only input - accepts ephemeral values
  input_wo         = ephemeral.tls_private_key.ssh.public_key_openssh
  input_wo_version = 1

  # Optional: triggers replacement when value changes
  triggers_replace = [some_resource.id]
}

# Output - readable from state
output "public_key" {
  value = wodata_wo_data.example.output
}
```

## Developing the Provider

To compile the provider, run `go install`.

To run tests:

```shell
go test ./...
```

To run acceptance tests (creates real resources):

```shell
make testacc
```

To generate documentation:

```shell
make generate
```
