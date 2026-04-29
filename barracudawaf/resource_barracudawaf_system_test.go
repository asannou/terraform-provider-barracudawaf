package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBarrcudaWAFSystem(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAcctPreCheck(t) },
		Providers:         testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBarrcudaWAFSystemConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("barracudawaf_system.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("barracudawaf_system.test", "hostname", "waf"),
					resource.TestCheckResourceAttr("barracudawaf_system.test", "locale", "English"),
				),
			},
		},
	})
}

func testAccBarrcudaWAFSystemConfig() string {
	return fmt.Sprintf(`
resource "barracudawaf_system" "test" {
    domain   = "example.com"
    hostname = "waf"
    locale   = "English"
}
`)
}
