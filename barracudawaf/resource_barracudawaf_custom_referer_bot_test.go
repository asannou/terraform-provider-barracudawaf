package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var CUSTOM_REFERER_BOT_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_custom_referer_bot" "test_bot" {
    name = "test_terraform_bot"
}
`

func TestAccBarracudaWAFCustomRefererBot_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CUSTOM_REFERER_BOT_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckCustomRefererBotExists("test_terraform_bot"),
					resource.TestCheckResourceAttr("barracudawaf_custom_referer_bot.test_bot", "name", "test_terraform_bot"),
				),
			},
		},
	})
}

func testCheckCustomRefererBotExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/bot-spam-types"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("custom referer bot %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("custom referer bot (%s) not found on the system", name)
		}

		return nil
	}
}
