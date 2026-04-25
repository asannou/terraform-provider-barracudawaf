package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var RESPONSE_PAGE_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_response_pages" "demo_page_1" {
    name        = "DemoPage1"
    status_code = "200"
    type        = "Other Pages"
    body        = "<html><body><h1>Demo Page</h1></body></html>"
}
`

func TestAccBarracudaWAFResponsePage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: RESPONSE_PAGE_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckResponsePageExists("DemoPage1"),
					resource.TestCheckResourceAttr("barracudawaf_response_pages.demo_page_1", "name", "DemoPage1"),
					resource.TestCheckResourceAttr("barracudawaf_response_pages.demo_page_1", "status_code", "200"),
					resource.TestCheckResourceAttr("barracudawaf_response_pages.demo_page_1", "type", "Other Pages"),
				),
			},
		},
	})
}

func testCheckResponsePageExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/response-pages"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(name, request)
		if err != nil {
			return err
		}

		if resources == nil {
			return fmt.Errorf("response page %s was not created.", name)
		}

		var dataItems map[string]interface{}
		for _, dataItems = range resources.Data {
			if dataItems["name"] == name {
				break
			}
		}

		if dataItems["name"] != name {
			return fmt.Errorf("response page (%s) not found on the system", name)
		}

		return nil
	}
}
