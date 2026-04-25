package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var GEO_POOL_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_geo_pool" "test_geo" {
    name = "test_terraform_geo"
    region = "Asia"
}
`

func TestAccBarracudaWAFGEOPool_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: GEO_POOL_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckGEOPoolExists("test_terraform_geo"),
					resource.TestCheckResourceAttr("barracudawaf_geo_pool.test_geo", "name", "test_terraform_geo"),
				),
			},
		},
	})
}

func testCheckGEOPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/geo-pools"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("geo pool %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("geo pool (%s) not found on the system", name)
		}

		return nil
	}
}
