package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBarracudaWAFSystemExportLogSettings_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBarracudaWAFSystemExportLogSettingsConfig("Enable"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSystemExportLogSettingsExists("barracudawaf_system_export_log_settings.test"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_settings.test", "export_access_logs", "Enable"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_settings.test", "export_audit_logs", "Enable"),
				),
			},
			{
				Config: testAccBarracudaWAFSystemExportLogSettingsConfig("Disable"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSystemExportLogSettingsExists("barracudawaf_system_export_log_settings.test"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_settings.test", "export_access_logs", "Disable"),
					resource.TestCheckResourceAttr("barracudawaf_system_export_log_settings.test", "export_audit_logs", "Disable"),
				),
			},
		},
	})
}

func testCheckSystemExportLogSettingsExists(n string) resource.TestCheckFunc {
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
			URL:    "/system/export-log-settings",
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

func testAccBarracudaWAFSystemExportLogSettingsConfig(status string) string {
	return fmt.Sprintf(`
resource "barracudawaf_system_export_log_settings" "test" {
  export_access_logs           = "%s"
  export_audit_logs            = "%s"
  export_web_firewall_logs     = "%s"
  export_network_firewall_logs = "%s"
  export_system_logs           = "%s"
}
`, status, status, status, status, status)
}
