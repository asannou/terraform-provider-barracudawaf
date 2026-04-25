package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var ADMIN_IP_RANGE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_admin_ip_range" "demo_admin_ip_range_1" {
    ip_address = "192.168.10.0"
    netmask    = "255.255.255.0"
}
`

func TestAccBarracudaWAFAdminIPRange_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ADMIN_IP_RANGE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckAdminIPRangeExists("192.168.10.0"),
					resource.TestCheckResourceAttr("barracudawaf_admin_ip_range.demo_admin_ip_range_1", "ip_address", "192.168.10.0"),
					resource.TestCheckResourceAttr("barracudawaf_admin_ip_range.demo_admin_ip_range_1", "netmask", "255.255.255.0"),
				),
			},
		},
	})
}

func testCheckAdminIPRangeExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/admin-ip-range"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("admin ip range %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["ip-address"] == name {
				break
			}
		}

		if dataItems["ip-address"] != name {
			return fmt.Errorf("admin ip range (%s) not found on the system", name)
		}

		return nil
	}
}
