package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var DATA_THEFT_PROTECTION_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_security_policies" "test_pol_for_dtp" {
    name = "PolForDTP"
}

resource "barracudawaf_data_theft_protection" "test_dtp" {
    name = "test_terraform_dtp"
    parent = [ barracudawaf_security_policies.test_pol_for_dtp.name ]
    enable = "Yes"
    identity_theft_type = "Credit Cards"
}
`

func TestAccBarracudaWAFDataTheftProtection_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DATA_THEFT_PROTECTION_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDataTheftProtectionExists("PolForDTP", "test_terraform_dtp"),
					resource.TestCheckResourceAttr("barracudawaf_data_theft_protection.test_dtp", "name", "test_terraform_dtp"),
				),
			},
		},
	})
}

func testCheckDataTheftProtectionExists(policyName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := fmt.Sprintf("/security-policies/%s/protected-data-types", policyName)
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("data theft protection %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("data theft protection (%s) not found on the system", name)
		}

		return nil
	}
}
