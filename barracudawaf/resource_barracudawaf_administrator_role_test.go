package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var ADMINISTRATOR_ROLE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_administrator_role" "test_role" {
    name = "test_terraform_role"
    api_privilege = "Yes"
    role_type = "Regular"
}
`

func TestAccBarracudaWAFAdministratorRole_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ADMINISTRATOR_ROLE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckAdministratorRoleExists("test_terraform_role"),
					resource.TestCheckResourceAttr("barracudawaf_administrator_role.test_role", "name", "test_terraform_role"),
					resource.TestCheckResourceAttr("barracudawaf_administrator_role.test_role", "api_privilege", "Yes"),
				),
			},
		},
	})
}

func testCheckAdministratorRoleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/administrator-roles"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("administrator role %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("administrator role (%s) not found on the system", name)
		}

		return nil
	}
}
