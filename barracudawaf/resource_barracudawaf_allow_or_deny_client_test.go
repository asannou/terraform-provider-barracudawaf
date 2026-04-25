package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var ALLOW_OR_DENY_CLIENT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "test_svc_for_adc" {
    name = "SvcForADC"
    ip_address = "172.30.1.14"
    port = "80"
}

resource "barracudawaf_allow_or_deny_client" "test_adc" {
    name = "test_terraform_adc"
    parent = [ barracudawaf_services.test_svc_for_adc.name ]
    action = "Deny"
}
`

func TestAccBarracudaWAFAllowOrDenyClient_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ALLOW_OR_DENY_CLIENT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckAllowOrDenyClientExists("SvcForADC", "test_terraform_adc"),
					resource.TestCheckResourceAttr("barracudawaf_allow_or_deny_client.test_adc", "name", "test_terraform_adc"),
				),
			},
		},
	})
}

func testCheckAllowOrDenyClientExists(serviceName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := fmt.Sprintf("/services/%s/allow-deny-clients", serviceName)
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("allow or deny client %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("allow or deny client (%s) not found on the system", name)
		}

		return nil
	}
}
