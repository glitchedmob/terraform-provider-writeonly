# Test Plan for writeonly_data Resource

## Test Categories

### 1. Basic Functionality Tests

#### TestAccWriteonlyData_Basic
Create with a plain string value, verify `output` captures it, verify `id` is generated.

#### TestAccWriteonlyData_OutputMatchesInput
Verify `output` attribute equals the `input_wo` value exactly.

---

### 2. Write-Only Verification (Critical)

#### TestAccWriteonlyData_WriteOnlyNotInState
Apply config, then use `terraform state pull` to verify `input_wo` is NOT in state (should be null).

#### TestAccWriteonlyData_OutputInState
Verify `output` IS persisted in state while `input_wo` is not.

---

### 3. Version Trigger Tests (Critical)

#### TestAccWriteonlyData_VersionTriggersUpdate
1. Create with `input_wo = "value1"`, `input_wo_version = 1`
2. Update to `input_wo = "value2"`, `input_wo_version = 2`
3. Verify `output` changes to "value2"

#### TestAccWriteonlyData_NoVersionChangeNoUpdate
Apply same config twice (no version change). Verify no detected change using `plancheck.ExpectNoChange`.

#### TestAccWriteonlyData_InputChangeWithoutVersionNoUpdate (Missing - User Requested)
1. Create with `input_wo = "value1"`, `input_wo_version = 1`
2. Update config to `input_wo = "value2"`, keep `input_wo_version = 1`
3. Verify `output` remains "value1" (unchanged)

This is critical because it proves that without incrementing the version, Terraform cannot detect changes to write-only values - they are simply ignored.

---

### 4. Ephemeral Resource Test (Primary Use Case)

#### TestAccWriteonlyData_WithEphemeral
```hcl
ephemeral "tls_private_key" "example" { 
  algorithm = "ED25519" 
}

resource "writeonly_data" "example" { 
  input_wo = ephemeral.tls_private_key.example.public_key_openssh
  input_wo_version = 1
}
```
Verify public key is captured in `output`.

---

### 5. triggers_replace Test

#### TestAccWriteonlyData_TriggersReplace
Similar to `terraform_data` - when `triggers_replace` changes, resource should be replaced (not updated). Use `ExpectResourceAction(plancheck.ResourceActionReplace)`.

---

### 6. Edge Cases

#### TestAccWriteonlyData_WithoutInputWO
What happens when `input_wo` is not provided? `output` should likely be null/empty.

#### TestAccWriteonlyData_UpdateOutputWithoutInputChange
Updating other fields (like `triggers_replace`) shouldn't re-capture `input_wo`.

#### TestAccWriteonlyData_EmptyStringInput
Test with empty string as input_wo - should be captured as empty in output.

---

## Test Execution

Run tests with:
```bash
TF_ACC=1 go test -v ./...
```

## Test Priority

1. **Must Have** (Core Value):
   - TestAccWriteonlyData_WriteOnlyNotInState
   - TestAccWriteonlyData_OutputInState
   - TestAccWriteonlyData_VersionTriggersUpdate
   - TestAccWriteonlyData_InputChangeWithoutVersionNoUpdate
   - TestAccWriteonlyData_WithEphemeral

2. **Should Have**:
   - TestAccWriteonlyData_Basic
   - TestAccWriteonlyData_TriggersReplace
   - TestAccWriteonlyData_NoVersionChangeNoUpdate

3. **Nice to Have**:
   - TestAccWriteonlyData_WithoutInputWO
   - TestAccWriteonlyData_UpdateOutputWithoutInputChange
   - TestAccWriteonlyData_EmptyStringInput

**Note**: Import is not supported, matching the behavior of `terraform_data`.
