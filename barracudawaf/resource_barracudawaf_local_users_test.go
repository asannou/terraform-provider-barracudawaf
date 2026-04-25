package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var LOCAL_USERS_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_local_users" "test_user" {
    name = "test_terraform_user"
    password = "Password123!"
}
`

func TestAccBarracudaWAFLocalUsers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: LOCAL_USERS_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckLocalUsersExists("test_terraform_user"),
					resource.TestCheckResourceAttr("barracudawaf_local_users.test_user", "name", "test_terraform_user"),
				),
			},
		},
	})
}

func testCheckLocalUsersExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/local-users"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("local user %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("local user (%s) not found on the system", name)
		}

		return nil
	}
}
