package barracudawaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var ATTACK_TYPE_PATTERN_RESOURCE_CREATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_custom_attack_type" "test_attack_type" {
    name = "test_custom_attack_for_pattern"
}

resource "barracudawaf_attack_type_pattern" "test_pattern" {
    name           = "test_pattern_1"
    regex          = ".*test.*"
    mode           = "Passive"
    description    = "Test Pattern"
    algorithm      = "None"
    case_sensitive = "No"
    parent         = [ barracudawaf_custom_attack_type.test_attack_type.name ]
}
`

var ATTACK_TYPE_PATTERN_RESOURCE_UPDATE = BARRACUDA_WAF_PROVIDER + `
resource "barracudawaf_custom_attack_type" "test_attack_type" {
    name = "test_custom_attack_for_pattern"
}

resource "barracudawaf_attack_type_pattern" "test_pattern" {
    name           = "test_pattern_1"
    regex          = ".*test_updated.*"
    mode           = "Active"
    description    = "Updated Test Pattern"
    algorithm      = "None"
    case_sensitive = "Yes"
    parent         = [ barracudawaf_custom_attack_type.test_attack_type.name ]
}
`

func TestAccBarracudaWAFAttackTypePattern_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAcctPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ATTACK_TYPE_PATTERN_RESOURCE_CREATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckAttackTypePatternExists("test_custom_attack_for_pattern", "test_pattern_1"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "name", "test_pattern_1"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "regex", ".*test.*"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "mode", "Passive"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "case_sensitive", "No"),
				),
			},
			{
				Config: ATTACK_TYPE_PATTERN_RESOURCE_UPDATE,
				Check: resource.ComposeTestCheckFunc(
					testCheckAttackTypePatternExists("test_custom_attack_for_pattern", "test_pattern_1"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "regex", ".*test_updated.*"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "mode", "Active"),
					resource.TestCheckResourceAttr("barracudawaf_attack_type_pattern.test_pattern", "case_sensitive", "Yes"),
				),
			},
		},
	})
}

func testCheckAttackTypePatternExists(attackTypeName string, patternName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*BarracudaWAF)

		resourceEndpoint := "/attack-types/" + attackTypeName + "/attack-patterns"
		request := &APIRequest{
			Method: "get",
			URL:    resourceEndpoint,
		}

		resources, err := client.GetBarracudaWAFResource(patternName, request)
		if err != nil {
			return err
		}

		if resources == nil || resources.Data == nil {
			return fmt.Errorf("attack type pattern %s was not created.", patternName)
		}

		found := false
		for _, dataItems := range resources.Data {
			if dataItems["name"] == patternName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("attack type pattern (%s) not found on the system", patternName)
		}

		return nil
	}
}
