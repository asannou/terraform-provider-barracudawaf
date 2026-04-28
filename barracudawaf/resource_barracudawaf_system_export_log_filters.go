package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFSystemExportLogFilters() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSystemExportLogFiltersUpdate,
		Read:   resourceCudaWAFSystemExportLogFiltersRead,
		Update: resourceCudaWAFSystemExportLogFiltersUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"web_firewall_log_severity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "7-Debug",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validSeverities := []string{
						"0-Emergency", "1-Alert", "2-Critical", "3-Error",
						"4-Warning", "5-Notice", "6-Information", "7-Debug",
					}
					found := false
					for _, s := range validSeverities {
						if v == s {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validSeverities))
					}
					return
				},
				Description: "Web Firewall Log Severity",
			},
			"system_log_severity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "7-Debug",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validSeverities := []string{
						"0-Emergency", "1-Alert", "2-Critical", "3-Error",
						"4-Warning", "5-Notice", "6-Information", "7-Debug",
					}
					found := false
					for _, s := range validSeverities {
						if v == s {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validSeverities))
					}
					return
				},
				Description: "System Log Severity",
			},
		},

		Description: "`barracudawaf_system_export_log_filters` manages `System Export Log Filters` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSystemExportLogFiltersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/system/export-log-filters"
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
	dataItems = resources.Data["export-log-filters"]
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

	d.Set("web_firewall_log_severity", dataItems["web-firewall-log-severity"])
	d.Set("system_log_severity", dataItems["system-log-severity"])

	return nil
}

func resourceCudaWAFSystemExportLogFiltersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	// For singleton resources, we can use a fixed ID or any string.
	name := "system_export_log_filters"
	if d.Id() == "" {
		d.SetId(name)
	}

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/system/export-log-filters"
	request := hydrateBarracudaWAFSystemExportLogFiltersResource(d, "put", resourceEndpoint)
	_, err := client.putReq(request.Body, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSystemExportLogFiltersRead(d, m)
}

func hydrateBarracudaWAFSystemExportLogFiltersResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"web-firewall-log-severity": d.Get("web_firewall_log_severity").(string),
		"system-log-severity":       d.Get("system_log_severity").(string),
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
