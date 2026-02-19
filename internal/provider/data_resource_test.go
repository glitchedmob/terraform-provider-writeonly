// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataResourceConfig("test-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("writeonly_data.test", "input_wo", "test-value"),
					resource.TestCheckResourceAttr("writeonly_data.test", "output", "test-value"),
					resource.TestCheckResourceAttrSet("writeonly_data.test", "id"),
				),
			},
			{
				ResourceName:            "writeonly_data.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"input_wo", "output"},
			},
		},
	})
}

func testAccDataResourceConfig(value string) string {
	return `
resource "writeonly_data" "test" {
  input_wo = "` + value + `"
}
`
}
