package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var CSP_POLICY_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "test_svc_for_csp" {
    name = "SvcForCSP"
    ip_address = "172.30.1.13"
    port = "80"
}

resource "barracudawaf_csp_policy" "test_csp" {
    csp_policy_name = "test_terraform_csp"
    parent = [ barracudawaf_services.test_svc_for_csp.name, "default" ]
    status = "On"
}
`

func TestAccBarracudaWAFCSPPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CSP_POLICY_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckCSPPolicyExists("SvcForCSP", "default", "test_terraform_csp"),
					resource.TestCheckResourceAttr("barracudawaf_csp_policy.test_csp", "csp_policy_name", "test_terraform_csp"),
				),
			},
		},
	})
}

func testCheckCSPPolicyExists(serviceName string, ruleName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := fmt.Sprintf("/services/%s/client-side-protection/%s/csp-policies", serviceName, ruleName)
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("csp policy %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["csp-policy-name"] == name {
				break
			}
		}

		if dataItems["csp-policy-name"] != name {
			return fmt.Errorf("csp policy (%s) not found on the system", name)
		}

		return nil
	}
}
