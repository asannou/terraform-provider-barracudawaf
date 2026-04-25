package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var JSON_SECURITY_POLICY_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_json_security_policy" "test_json_policy" {
    name = "test_terraform_json_policy"
    max_keys = "100"
    max_key_length = "256"
    max_value_length = "1024"
}
`

func TestAccBarracudaWAFJSONSecurityPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: JSON_SECURITY_POLICY_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckJSONSecurityPolicyExists("test_terraform_json_policy"),
					resource.TestCheckResourceAttr("barracudawaf_json_security_policy.test_json_policy", "name", "test_terraform_json_policy"),
				),
			},
		},
	})
}

func testCheckJSONSecurityPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/json-security-policies"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("json security policy %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("json security policy (%s) not found on the system", name)
		}

		return nil
	}
}
