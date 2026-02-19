// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWoDataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWoDataResourceConfig("test-value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wodata_wo_data.test", "input_wo", "test-value"),
					resource.TestCheckResourceAttr("wodata_wo_data.test", "output", "test-value"),
					resource.TestCheckResourceAttrSet("wodata_wo_data.test", "id"),
				),
			},
			{
				ResourceName:            "wodata_wo_data.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"input_wo"},
			},
		},
	})
}

func testAccWoDataResourceConfig(value string) string {
	return `
resource "wodata_wo_data" "test" {
	input_wo = "` + value + `"
}
`
}
