package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWriteonlyData_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfig("test-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"writeonly_data.test", "id", regexp.MustCompile(`^writeonly_data_`)),
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "test-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_OutputMatchesInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfig("exact-match-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "exact-match-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_WriteOnlyNotInState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfig("secret-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckWriteOnlyNotInState("writeonly_data.test"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_OutputInState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfig("captured-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOutputInState("writeonly_data.test", "captured-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_VersionTriggersUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithVersion("value1", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "value1"),
				),
			},
			{
				Config: testAccWriteonlyDataConfigWithVersion("value2", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "value2"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_NoVersionChangeNoUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithVersion("stable-value", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "stable-value"),
				),
			},
			{
				Config: testAccWriteonlyDataConfigWithVersion("stable-value", 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "stable-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_InputChangeWithoutVersionNoUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithVersion("original-value", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "original-value"),
				),
			},
			{
				Config: testAccWriteonlyDataConfigWithVersion("changed-value", 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "original-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_WithEphemeral(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithEphemeral(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"writeonly_data.test", "output",
						regexp.MustCompile(`^ssh-ed25519\s+`)),
				),
			},
		},
	})
}

func TestAccWriteonlyData_TriggersReplace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithTriggersReplace("trigger-value", "a"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "trigger-value"),
				),
			},
			{
				Config: testAccWriteonlyDataConfigWithTriggersReplace("trigger-value", "b"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("writeonly_data.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "trigger-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_WithoutInputWO(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithoutInput(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"writeonly_data.test", "id", regexp.MustCompile(`^writeonly_data_`)),
					resource.TestCheckNoResourceAttr("writeonly_data.test", "output"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_UpdateOutputWithoutInputChange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfigWithTriggersReplace("stable-value", "initial"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "stable-value"),
				),
			},
			{
				Config: testAccWriteonlyDataConfigWithTriggersReplace("stable-value", "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "stable-value"),
				),
			},
		},
	})
}

func TestAccWriteonlyData_EmptyStringInput(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWriteonlyDataConfig(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "output", ""),
				),
			},
		},
	})
}

func testAccCheckWriteOnlyNotInState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		stateJSON, err := json.Marshal(rs.Primary.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal state: %w", err)
		}

		var state map[string]interface{}
		if err := json.Unmarshal(stateJSON, &state); err != nil {
			return fmt.Errorf("failed to unmarshal state: %w", err)
		}

		if v, exists := state["input_wo"]; exists && v != "" {
			return fmt.Errorf("input_wo should not be in state, got: %v", v)
		}

		return nil
	}
}

func testAccCheckOutputInState(resourceName, expectedValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		actualValue := rs.Primary.Attributes["output"]
		if actualValue != expectedValue {
			return fmt.Errorf("output should be %q in state, got: %q", expectedValue, actualValue)
		}

		return nil
	}
}

func testAccWriteonlyDataConfig(value string) string {
	return fmt.Sprintf(`
resource "writeonly_data" "test" {
  input_wo = %[1]q
}
`, value)
}

func testAccWriteonlyDataConfigWithVersion(value string, version int) string {
	return fmt.Sprintf(`
resource "writeonly_data" "test" {
  input_wo          = %[1]q
  input_wo_version = %[2]d
}
`, value, version)
}

func testAccWriteonlyDataConfigWithTriggersReplace(value, trigger string) string {
	return fmt.Sprintf(`
resource "terraform_data" "trigger" {
  input = %[2]q
}

resource "writeonly_data" "test" {
  input_wo          = %[1]q
  triggers_replace = [terraform_data.trigger.output]
}
`, value, trigger)
}

func testAccWriteonlyDataConfigWithoutInput() string {
	return `
resource "writeonly_data" "test" {
}
`
}

func testAccWriteonlyDataConfigWithEphemeral() string {
	return `
ephemeral "tls_private_key" "example" {
  algorithm = "ED25519"
}

resource "writeonly_data" "test" {
  input_wo          = ephemeral.tls_private_key.example.public_key_openssh
  input_wo_version = 1
}
`
}
