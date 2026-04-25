package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var VSITE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_vsites" "demo_vsite_1" {
    name      = "DemoVsite1"
    interface = "WAN"
    comments  = "Demo Vsite with Terraform"
}
`

func TestAccBarracudaWAFVsite_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: VSITE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckVsiteExists("DemoVsite1"),
					resource.TestCheckResourceAttr("barracudawaf_vsites.demo_vsite_1", "name", "DemoVsite1"),
					resource.TestCheckResourceAttr("barracudawaf_vsites.demo_vsite_1", "interface", "WAN"),
					resource.TestCheckResourceAttr("barracudawaf_vsites.demo_vsite_1", "comments", "Demo Vsite with Terraform"),
				),
			},
		},
	})
}

func testCheckVsiteExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/vsites"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil {
			return fmt.Errorf("vsite %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("vsite (%s) not found on the system", name)
		}

		return nil
	}
}
