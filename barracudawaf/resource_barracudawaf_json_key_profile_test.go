package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var JSON_KEY_PROFILE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_services" "test_svc_for_json_key" {
    name = "SvcForJSONKey"
    ip_address = "172.30.1.11"
    port = "80"
}

resource "barracudawaf_json_profile" "test_json_profile_for_key" {
    name = "test_json_profile"
    parent = [ barracudawaf_services.test_svc_for_json_key.name ]
    host_match = "example.com"
    url_match = "/api"
    method = [ "POST" ]
}

resource "barracudawaf_json_key_profile" "test_json_key" {
    name = "test_terraform_json_key"
    parent = [ barracudawaf_services.test_svc_for_json_key.name, barracudawaf_json_profile.test_json_profile_for_key.name ]
    key = "user_id"
    value_class = "Integer"
}
`

func TestAccBarracudaWAFJSONKeyProfile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: JSON_KEY_PROFILE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckJSONKeyProfileExists("SvcForJSONKey", "test_json_profile", "test_terraform_json_key"),
					resource.TestCheckResourceAttr("barracudawaf_json_key_profile.test_json_key", "name", "test_terraform_json_key"),
				),
			},
		},
	})
}

func testCheckJSONKeyProfileExists(serviceName string, jsonProfileName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := fmt.Sprintf("/services/%s/json-profiles/%s/json-key-profiles", serviceName, jsonProfileName)
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("json key profile %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("json key profile (%s) not found on the system", name)
		}

		return nil
	}
}
