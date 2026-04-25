package barracudawaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var CUSTOM_IP_BLOCKLIST_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_custom_ip_blocklist" "test_blocklist" {
    custom_ip_list = "Download"
    download_url = "http://example.com/ips.txt"
}
`

func TestAccBarracudaWAFCustomIPBlocklist_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CUSTOM_IP_BLOCKLIST_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("barracudawaf_custom_ip_blocklist.test_blocklist", "custom_ip_list", "Download"),
					resource.TestCheckResourceAttr("barracudawaf_custom_ip_blocklist.test_blocklist", "download_url", "http://example.com/ips.txt"),
				),
			},
		},
	})
}
