package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var CUSTOM_ATTACK_TYPE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_custom_attack_type" "test_attack_type" {
    name = "test_custom_attack"
}
`

func TestAccBarracudaWAFCustomAttackType_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CUSTOM_ATTACK_TYPE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckCustomAttackTypeExists("test_custom_attack"),
					resource.TestCheckResourceAttr("barracudawaf_custom_attack_type.test_attack_type", "name", "test_custom_attack"),
				),
			},
		},
	})
}

func testCheckCustomAttackTypeExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/attack-types"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("custom attack type %s was not created.", name)
		}

		found := false
		for _, dataItems := range resources.Data {
			if dataItems["name"] == name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("custom attack type (%s) not found on the system", name)
		}

		return nil
	}
}
