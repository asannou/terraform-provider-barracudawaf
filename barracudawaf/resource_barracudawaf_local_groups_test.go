package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var LOCAL_GROUPS_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_local_groups" "test_group" {
    name = "test_terraform_group"
}
`

func TestAccBarracudaWAFLocalGroups_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: LOCAL_GROUPS_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckLocalGroupsExists("test_terraform_group"),
					resource.TestCheckResourceAttr("barracudawaf_local_groups.test_group", "name", "test_terraform_group"),
				),
			},
		},
	})
}

func testCheckLocalGroupsExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/local-groups"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("local group %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("local group (%s) not found on the system", name)
		}

		return nil
	}
}
