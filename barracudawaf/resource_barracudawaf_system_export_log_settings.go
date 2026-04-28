package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFSystemExportLogSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSystemExportLogSettingsUpdate,
		Read:   resourceCudaWAFSystemExportLogSettingsRead,
		Update: resourceCudaWAFSystemExportLogSettingsUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"export_access_logs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Disable",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Enable" && v != "Disable" {
						errs = append(errs, fmt.Errorf("%q must be either Enable or Disable", key))
					}
					return
				},
				Description: "Enable export of Access Logs",
			},
			"export_audit_logs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Disable",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Enable" && v != "Disable" {
						errs = append(errs, fmt.Errorf("%q must be either Enable or Disable", key))
					}
					return
				},
				Description: "Enable export of Audit Logs",
			},
			"export_web_firewall_logs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Disable",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Enable" && v != "Disable" {
						errs = append(errs, fmt.Errorf("%q must be either Enable or Disable", key))
					}
					return
				},
				Description: "Enable export of Web Firewall Logs",
			},
			"export_network_firewall_logs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Disable",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Enable" && v != "Disable" {
						errs = append(errs, fmt.Errorf("%q must be either Enable or Disable", key))
					}
					return
				},
				Description: "Enable export of Network Firewall Logs",
			},
			"export_system_logs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Disable",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Enable" && v != "Disable" {
						errs = append(errs, fmt.Errorf("%q must be either Enable or Disable", key))
					}
					return
				},
				Description: "Enable export of System Logs",
			},
		},

		Description: "`barracudawaf_system_export_log_settings` manages `System Export Log Settings` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSystemExportLogSettingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/system/export-log-settings"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	resources, err := client.GetBarracudaWAFResource(name, request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var dataItems map[string]interface{}
	dataItems = resources.Data["export-log-settings"]
	if dataItems == nil {
		// Fallback for different API structure if needed
		for _, v := range resources.Data {
			dataItems = v
			break
		}
	}

	if dataItems == nil {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("export_access_logs", dataItems["export-access-logs"])
	d.Set("export_audit_logs", dataItems["export-audit-logs"])
	d.Set("export_web_firewall_logs", dataItems["export-web-firewall-logs"])
	d.Set("export_network_firewall_logs", dataItems["export-network-firewall-logs"])
	d.Set("export_system_logs", dataItems["export-system-logs"])

	return nil
}

func resourceCudaWAFSystemExportLogSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	// For singleton resources, we can use a fixed ID or any string.
	// We'll use "system_export_log_settings" as the ID if it's not set.
	name := "system_export_log_settings"
	if d.Id() == "" {
		d.SetId(name)
	}

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/system/export-log-settings"
	request := hydrateBarracudaWAFSystemExportLogSettingsResource(d, "put", resourceEndpoint)
	_, err := client.putReq(request.Body, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSystemExportLogSettingsRead(d, m)
}

func hydrateBarracudaWAFSystemExportLogSettingsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"export-access-logs":           d.Get("export_access_logs").(string),
		"export-audit-logs":            d.Get("export_audit_logs").(string),
		"export-web-firewall-logs":     d.Get("export_web_firewall_logs").(string),
		"export-network-firewall-logs": d.Get("export_network_firewall_logs").(string),
		"export-system-logs":           d.Get("export_system_logs").(string),
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if len(val) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
