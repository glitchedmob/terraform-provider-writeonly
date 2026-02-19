resource "writeonly_data" "example" {
  # Can accept ephemeral values:
  # input_wo = ephemeral.tls_private_key.ssh.public_key_openssh
  input_wo         = "example-write-only-value"
  input_wo_version = 1
}
