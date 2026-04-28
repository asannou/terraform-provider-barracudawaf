package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBarrcudaWAFWebScrapingPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWebScrapingPolicyConfig("test_scraping_policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("barracudawaf_web_scraping_policy.test_scraping_policy", "name", "test_scraping_policy"),
					resource.TestCheckResourceAttr("barracudawaf_web_scraping_policy.test_scraping_policy", "comments", "Terraform Managed Web Scraping Policy"),
				),
			},
		},
	})
}

func testAccWebScrapingPolicyConfig(name string) string {
	return fmt.Sprintf(`
resource "barracudawaf_web_scraping_policy" "%s" {
  name                          = "%s"
  comments                      = "Terraform Managed Web Scraping Policy"
  detect_mouse_event            = "Yes"
  insert_delay                  = "Yes"
  whitelisted_bots              = ["Google", "Baidu"]
  blacklisted_categories        = ["Web scraper", "Vulnerability scanner"]
  delay_time                    = "10"
  insert_disallowed_urls        = "No"
  insert_hidden_links           = "Yes"
  insert_javascript_in_response = "Yes"
}
`, name, name)
}
