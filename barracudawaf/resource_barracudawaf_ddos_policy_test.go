package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var DDOS_POLICY_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "test_svc_for_ddos" {
    name = "SvcForDDoS"
    ip_address = "172.30.1.12"
    port = "80"
}

resource "barracudawaf_ddos_policy" "test_ddos" {
    name = "test_terraform_ddos"
    parent = [ barracudawaf_services.test_svc_for_ddos.name ]
    evaluate_clients = "On"
}
`

func TestAccBarracudaWAFDDoSPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: DDOS_POLICY_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckDDoSPolicyExists("SvcForDDoS", "test_terraform_ddos"),
					resource.TestCheckResourceAttr("barracudawaf_ddos_policy.test_ddos", "name", "test_terraform_ddos"),
				),
			},
		},
	})
}

func testCheckDDoSPolicyExists(serviceName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := fmt.Sprintf("/services/%s/ddos-policies", serviceName)
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("ddos policy %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("ddos policy (%s) not found on the system", name)
		}

		return nil
	}
}
