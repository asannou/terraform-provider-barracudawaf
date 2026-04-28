package barracudawaf

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider : Schema definition for barracudawaf provider
func Provider() *schema.Provider {

	// The actual provider
	provider := &schema.Provider{

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP Address of the WAF to be configured",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Admin port on the WAF to be configured",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the WAF to be configured",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password of the WAF to be configured",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"barracudawaf_trusted_ca_certificate":     resourceCudaWAFTrustedCaCertificate(),
			"barracudawaf_content_rules":              resourceCudaWAFContentRules(),
			"barracudawaf_trusted_server_certificate": resourceCudaWAFTrustedServerCertificate(),
			"barracudawaf_services":                   resourceCudaWAFServices(),
			"barracudawaf_content_rule_servers":       resourceCudaWAFContentRuleServers(),
			"barracudawaf_security_policies":          resourceCudaWAFSecurityPolicies(),
			"barracudawaf_signed_certificate":         resourceCudaWAFSignedCertificate(),
			"barracudawaf_self_signed_certificate":    resourceCudaWAFSelfSignedCertificate(),
			"barracudawaf_servers":                    resourceCudaWAFServers(),
			"barracudawaf_letsencrypt_certificate":    resourceCudaWAFLetsEncryptCertificate(),
			"barracudawaf_admin_ip_range":             resourceCudaWAFAdminIPRange(),
			"barracudawaf_response_pages":             resourceCudaWAFResponsePages(),
			"barracudawaf_vsites":                     resourceCudaWAFVsites(),
			"barracudawaf_administrator_role":         resourceCudaWAFAdministratorRole(),
			"barracudawaf_local_users":                resourceCudaWAFLocalUsers(),
			"barracudawaf_local_groups":               resourceCudaWAFLocalGroups(),
			"barracudawaf_json_security_policy":       resourceCudaWAFJSONSecurityPolicy(),
			"barracudawaf_json_profile":               resourceCudaWAFJSONProfile(),
			"barracudawaf_json_key_profile":           resourceCudaWAFJSONKeyProfile(),
			"barracudawaf_ddos_policy":                resourceCudaWAFDDoSPolicy(),
			"barracudawaf_csp_policy":                 resourceCudaWAFCSPPolicy(),
			"barracudawaf_data_theft_protection":      resourceCudaWAFDataTheftProtection(),
			"barracudawaf_custom_referer_bot":         resourceCudaWAFCustomRefererBot(),
			"barracudawaf_custom_ip_blocklist":        resourceCudaWAFCustomIPBlocklist(),
			"barracudawaf_allow_or_deny_client":       resourceCudaWAFAllowOrDenyClient(),
			"barracudawaf_geo_pool":                   resourceCudaWAFGEOPool(),
			"barracudawaf_http_request_rewrite":       resourceCudaWAFHTTPRequestRewrite(),
			"barracudawaf_http_response_rewrite":      resourceCudaWAFHTTPResponseRewrite(),
			"barracudawaf_response_body_rewrite":      resourceCudaWAFResponseBodyRewrite(),
			"barracudawaf_rate_control_pool":          resourceCudaWAFRateControlPool(),
			"barracudawaf_network_interface":          resourceCudaWAFNetworkInterface(),
			"barracudawaf_network_vlan":               resourceCudaWAFNetworkVLAN(),
			"barracudawaf_bond":                       resourceCudaWAFBond(),
			"barracudawaf_interface_route":            resourceCudaWAFInterfaceRoute(),
			"barracudawaf_destination_nat":            resourceCudaWAFDestinationNAT(),
			"barracudawaf_cluster":                    resourceCudaWAFCluster(),
			"barracudawaf_cluster_nodes":              resourceCudaWAFClusterNodes(),
			"barracudawaf_ntp_server":                 resourceCudaWAFNTPServer(),
			"barracudawaf_network_acl":                resourceCudaWAFNetworkACL(),
			"barracudawaf_adaptive_profiling_rule":    resourceCudaWAFAdaptiveProfilingRule(),
			"barracudawaf_parameter_optimizer":        resourceCudaWAFParameterOptimizer(),
			"barracudawaf_access_rule":                resourceCudaWAFAccessRule(),
			"barracudawaf_header_acl":                 resourceCudaWAFHeaderACL(),
			"barracudawaf_form_spam":                  resourceCudaWAFFormSpam(),
			"barracudawaf_web_scraping_policy":        resourceCudaWAFWebScrapingPolicy(),
			"barracudawaf_system_export_log_settings": resourceCudaWAFSystemExportLogSettings(),
			"barracudawaf_system_export_log_filters":  resourceCudaWAFSystemExportLogFilters(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		IPAddress: d.Get("address").(string),
		AdminPort: d.Get("port").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}
	cfg, err := config.Client()
	if err != nil {
		return cfg, err
	}
	cfg.UserAgent = fmt.Sprintf("Terraform/%s", terraformVersion)
	return cfg, err
}
