package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var JSON_PROFILE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "test_svc_for_json" {
    name = "SvcForJSON"
    ip_address = "172.30.1.10"
    port = "80"
}

resource "barracudawaf_json_profile" "test_json_profile" {
    name = "test_terraform_json_profile"
    parent = [ barracudawaf_services.test_svc_for_json.name ]
    host_match = "example.com"
    url_match = "/api/v1"
    method = [ "POST", "PUT" ]
}
`

func TestAccBarracudaWAFJSONProfile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: JSON_PROFILE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckJSONProfileExists("SvcForJSON", "test_terraform_json_profile"),
					resource.TestCheckResourceAttr("barracudawaf_json_profile.test_json_profile", "name", "test_terraform_json_profile"),
				),
			},
		},
	})
}

func testCheckJSONProfileExists(serviceName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := fmt.Sprintf("/services/%s/json-profiles", serviceName)
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("json profile %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("json profile (%s) not found on the system", name)
		}

		return nil
	}
}
