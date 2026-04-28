package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBarracudaWAFSystemExportLogFilters_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBarracudaWAFSystemExportLogFiltersConfig("5-Notice"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSystemExportLogFiltersExists("barracudawaf_system_export_log_filters.test"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_filters.test", "web_firewall_log_severity", "5-Notice"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_filters.test", "system_log_severity", "5-Notice"),
				),
			},
			{
				Config: testAccBarracudaWAFSystemExportLogFiltersConfig("7-Debug"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSystemExportLogFiltersExists("barracudawaf_system_export_log_filters.test"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_filters.test", "web_firewall_log_severity", "7-Debug"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_filters.test", "system_log_severity", "7-Debug"),
				),
			},
		},
	})
}

func testCheckSystemExportLogFiltersExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*BarracudaWAF)
		request := &APIRequest{
			Method: "get",
			URL:    "/system/export-log-filters",
		}

		resources, err := client.GetBarracudaWAFResource(rs.Primary.ID, request)
		if err != nil {
			return err
		}

		if resources.Data == nil {
			return fmt.Errorf("Resource not found")
		}

		return nil
	}
}

func testAccBarracudaWAFSystemExportLogFiltersConfig(severity string) string {
	return fmt.Sprintf(`
resource "barracudawaf_system_export_log_filters" "test" {
  web_firewall_log_severity = "%s"
  system_log_severity       = "%s"
}
`, severity, severity)
}
